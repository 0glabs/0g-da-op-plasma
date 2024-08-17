package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	pb "github.com/0glabs/0g-da-client/api/grpc/disperser"
	plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const VersionByte = 0x47
const RequestTimeout = 180 * time.Second
const MaxReties = 90

type ZgConfig struct {
	server string
}

type ZgStore struct {
	cfg ZgConfig
	log log.Logger
}

func NewZgStore(ctx context.Context, cfg ZgConfig, log log.Logger) (*ZgStore, error) {
	return &ZgStore{
		cfg: cfg,
		log: log,
	}, nil
}

func (s *ZgStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	epoch := binary.LittleEndian.Uint64(key[2:10])
	quorumId := binary.LittleEndian.Uint64(key[10:18])
	commit := key[10:]

	ctxWithTimeout, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()
	conn, err := grpc.DialContext(
		ctxWithTimeout,
		s.cfg.server,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), // 1 GiB
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial encoder: %w", err)
	}
	defer conn.Close()

	client := pb.NewDisperserClient(conn)
	result, err := client.RetrieveBlob(ctx, &pb.RetrieveBlobRequest{
		StorageRoot: commit,
		Epoch:       epoch,
		QuorumId:    quorumId,
	})
	if err != nil {
		return nil, err
	}

	return result.GetData(), nil
}

func (s *ZgStore) Put(ctx context.Context, value []byte) ([]byte, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()
	conn, err := grpc.DialContext(
		ctxWithTimeout,
		s.cfg.server,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), // 1 GiB
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial encoder: %w", err)
	}
	defer conn.Close()

	client := pb.NewDisperserClient(conn)

	reply, err := client.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:           value,
		SecurityParams: []*pb.SecurityParams{},
		TargetRowNum:   0,
	})
	if err != nil {
		return nil, err
	}

	return s.WaitBlobConfirmed(ctx, reply.GetRequestId())
}

func (s *ZgStore) WaitBlobConfirmed(ctx context.Context, requestId []byte) ([]byte, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, RequestTimeout*2)
	defer cancel()
	conn, err := grpc.DialContext(
		ctxWithTimeout,
		s.cfg.server,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), // 1 GiB
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial encoder: %w", err)
	}
	defer conn.Close()

	client := pb.NewDisperserClient(conn)

	var retryCount uint64
	for {
		blobStatus, err := client.GetBlobStatus(ctx, &pb.BlobStatusRequest{
			RequestId: requestId,
		})

		if err != nil {
			retryCount++
			if retryCount > MaxReties {
				return nil, fmt.Errorf("failed to get blob status: %w", err)
			}

			log.Error("failed to get blob status", "err", err)
			time.Sleep(3 * time.Second)
		}

		status := blobStatus.GetStatus()
		if status == pb.BlobStatus_CONFIRMED || status == pb.BlobStatus_FINALIZED {
			blobHeader := blobStatus.GetInfo().GetBlobHeader()

			var comm []byte
			epoch := make([]byte, 8)
			binary.LittleEndian.PutUint64(epoch, blobHeader.GetEpoch())
			comm = append(comm, epoch...)

			quorumId := make([]byte, 8)
			binary.LittleEndian.PutUint64(quorumId, blobHeader.GetQuorumId())
			comm = append(comm, quorumId...)

			comm = append(comm, blobHeader.GetDataRoot()...)
			commitment := plasma.NewGenericCommitment(append([]byte{VersionByte}, comm...))
			return commitment.Encode(), nil
		}

		if status == pb.BlobStatus_FAILED {
			return nil, fmt.Errorf("failed to put blob")
		}

		retryCount++
		if retryCount > MaxReties {
			return nil, fmt.Errorf("failed to get blob status, retry reached")
		}

		time.Sleep(3 * time.Second)
	}
}
