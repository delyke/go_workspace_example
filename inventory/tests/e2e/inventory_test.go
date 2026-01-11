//go:build integration
// +build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

var _ = Describe("Inventory Service", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).NotTo(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// чистим коллекцию после теста
		// err := env.ClearPartsCollection(ctx)
		// Expect(err).NotTo(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("Get", func() {
		It("должен успешно получать созданную при инициализации деталь", func() {
			partUuid := "55555555-5555-5555-5555-555555555555"

			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUuid,
			})

			Expect(err).ToNot(HaveOccurred())
			part := resp.GetPart()
			Expect(part).ToNot(BeNil())
			Expect(part.GetUuid()).To(Equal("55555555-5555-5555-5555-555555555555"))
			Expect(part.GetName()).To(Equal("Wing Module R"))
			Expect(part.GetCreatedAt()).ToNot(BeNil())
		})

		It("должен выдавать ошибку, так как детали не существует", func() {
			partUuid := "88888888-7777-8888-8888-888888888888"
			_, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUuid,
			})
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.NotFound))
		})
	})
})
