package user

func (s *Service) GetPayoutChannel(userID string) (*UserPayoutChannel, error) {
	var userPayoutChannel *UserPayoutChannel
	err := s.DB.Take(&userPayoutChannel, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return userPayoutChannel, nil
}

func (s *Service) SetPayoutChannel(userID string, userPayoutChannel *UserPayoutChannel) (*UserPayoutChannel, error) {
	userPayoutChannel.UserID = userID
	err := s.Validate.Struct(userPayoutChannel)
	if err != nil {
		return nil, err
	}
	err = s.DB.Save(&userPayoutChannel).Error
	return userPayoutChannel, err
}
