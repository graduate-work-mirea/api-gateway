package grpc

import (
    "context"
    "log"

    pb "api-gateway/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
    conn   *grpc.ClientConn
    client pb.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
    conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }
    client := pb.NewAuthServiceClient(conn)
    return &AuthClient{conn: conn, client: client}, nil
}

func (c *AuthClient) Close() {
    if err := c.conn.Close(); err != nil {
        log.Printf("Failed to close gRPC connection: %v", err)
    }
}

func (c *AuthClient) ValidateToken(token string) (bool, error) {
    resp, err := c.client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
    if err != nil {
        return false, err
    }
    return resp.Valid, nil
}
