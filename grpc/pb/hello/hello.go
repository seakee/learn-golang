package hello

import "context"

type Service struct {
	UnimplementedHelloServiceServer
}

// func (s *Service) mustEmbedUnimplementedHelloServiceServer() {
// 	panic("implement me")
// }

func (s *Service) Hello(ctx context.Context, args *String) (*String, error) {
	reply := &String{Value: "hello:" + args.GetValue()}

	return reply, nil
}
