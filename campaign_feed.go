package gadwords

type CampaignFeedService struct {
	Auth
}

func NewCampaignFeedService(auth *Auth) *CampaignFeedService {
	return &CampaignFeedService{Auth: *auth}
}
