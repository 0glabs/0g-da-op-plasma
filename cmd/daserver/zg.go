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
const SleepDuration = 3 * time.Second
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
	commit := key[18:]
	s.log.Debug("Get blob for", "epoch", epoch, "quorum", quorumId, "commit", commit)

	conn, err := grpc.NewClient(s.cfg.server, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024))) // 1 GiB
	if err != nil {
		return nil, fmt.Errorf("failed to dial 0g da client: %w", err)
	}
	defer conn.Close()

	client := pb.NewDisperserClient(conn)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()
	result, err := client.RetrieveBlob(ctxWithTimeout, &pb.RetrieveBlobRequest{
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
	conn, err := grpc.NewClient(s.cfg.server, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024))) // 1 GiB
	if err != nil {
		return nil, fmt.Errorf("failed to dial 0g da client: %w", err)
	}
	defer conn.Close()

	client := pb.NewDisperserClient(conn)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()
	reply, err := client.DisperseBlob(ctxWithTimeout, &pb.DisperseBlobRequest{
		Data: value,
	})
	if err != nil {
		return nil, err
	}

	return s.WaitBlobConfirmed(ctx, client, reply.GetRequestId())
}

func (s *ZgStore) WaitBlobConfirmed(ctx context.Context, client pb.DisperserClient, requestId []byte) ([]byte, error) {
	var retryCount uint64
	for {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, RequestTimeout)
		defer cancel()
		blobStatus, err := client.GetBlobStatus(ctxWithTimeout, &pb.BlobStatusRequest{
			RequestId: requestId,
		})

		if err != nil {
			retryCount++
			if retryCount > MaxReties {
				return nil, fmt.Errorf("failed to get blob status: %w", err)
			}

			log.Error("failed to get blob status", "err", err)
			time.Sleep(SleepDuration)
		}

		status := blobStatus.GetStatus()
		if status == pb.BlobStatus_FINALIZED {
			blobHeader := blobStatus.GetInfo().GetBlobHeader()

			var comm []byte
			epoch := make([]byte, 8)
			binary.LittleEndian.PutUint64(epoch, blobHeader.GetEpoch())
			comm = append(comm, epoch...)

			quorumId := make([]byte, 8)
			binary.LittleEndian.PutUint64(quorumId, blobHeader.GetQuorumId())
			comm = append(comm, quorumId...)

			comm = append(comm, blobHeader.GetStorageRoot()...)
			commitment := plasma.NewGenericCommitment(append([]byte{VersionByte}, comm...))
			return commitment.Encode(), nil
		}

		if status == pb.BlobStatus_FAILED {
			return nil, fmt.Errorf("failed to put blob")
		}

		retryCount++
		if retryCount > MaxReties {
			return nil, fmt.Errorf("failed to get blob status, maximum retry reached")
		}

		time.Sleep(SleepDuration)
	}
}
