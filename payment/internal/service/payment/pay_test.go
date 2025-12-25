package payment

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderUUID     = s.faker.UUID()
		userUUID      = s.faker.UUID()
		paymentMethod = s.faker.CreditCardType()
	)
	txID, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().NotEmpty(txID)
}
