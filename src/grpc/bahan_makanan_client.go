package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "app/src/grpc/proto/bahan_makanan"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BahanMakananClient struct {
	client pb.BahanMakananServiceClient
	conn   *grpc.ClientConn
}

func NewBahanMakananClient(serverAddr string) (*BahanMakananClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	client := pb.NewBahanMakananServiceClient(conn)
	return &BahanMakananClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *BahanMakananClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing gRPC connection: %v", err)
		}
	}
}

func (c *BahanMakananClient) GetAllBahanMakanan(ctx context.Context) (*pb.ListBahanMakananResponse, error) {
	return c.client.GetAllBahanMakanan(ctx, &pb.Empty{})
}

func (c *BahanMakananClient) GetBahanMakananByKode(ctx context.Context, kode string) (*pb.BahanMakananResponse, error) {
	return c.client.GetBahanMakananByKode(ctx, &pb.GetBahanMakananRequest{Kode: kode})
}

func (c *BahanMakananClient) GetBahanMakananById(ctx context.Context, id uint32) (*pb.BahanMakananResponse, error) {
	return c.client.GetBahanMakananById(ctx, &pb.GetBahanMakananByIdRequest{Id: id})
}

func (c *BahanMakananClient) GetBahanMakananByMentahOlahan(ctx context.Context, mentahOlahan string) (*pb.ListBahanMakananResponse, error) {
	return c.client.GetBahanMakananByMentahOlahan(ctx, &pb.GetBahanMakananByMentahOlahanRequest{MentahOlahan: mentahOlahan})
}

func (c *BahanMakananClient) GetBahanMakananByKelompok(ctx context.Context, kelompokMakanan string) (*pb.ListBahanMakananResponse, error) {
	return c.client.GetBahanMakananByKelompok(ctx, &pb.GetBahanMakananByKelompokRequest{KelompokMakanan: kelompokMakanan})
}

func (c *BahanMakananClient) UpdateBahanMakanan(ctx context.Context, id uint32, bahanMakanan *pb.BahanMakanan) (*pb.BahanMakananResponse, error) {
	return c.client.UpdateBahanMakanan(ctx, &pb.UpdateBahanMakananRequest{
		Id:           id,
		BahanMakanan: bahanMakanan,
	})
}
