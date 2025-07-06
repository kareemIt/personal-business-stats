package model

type EtsyJob struct {
	Profit         int
	TrackingNumber string
}

func NewEtsyJob(Profit int, trackingNumber string) *EtsyJob {
	if Profit == 0 || trackingNumber == "" {
		panic("Price cost and tracking number cannot be empty")
	}
	return &EtsyJob{
		Profit:         Profit,
		TrackingNumber: trackingNumber,
	}
}

func (p *EtsyJob) GetProfit() int {
	return p.Profit
}

func (p *EtsyJob) GetTrackingNumber() string {
	return p.TrackingNumber
}

func (p *EtsyJob) SetProfit(Profit int) {
	if Profit == 0 {
		panic("Price cost cannot be empty")
	}
	p.Profit = Profit
}
func (p *EtsyJob) SetTrackingNumber(trackingNumber string) {
	if trackingNumber == "" {
		panic("Tracking number cannot be empty")
	}
	p.TrackingNumber = trackingNumber
}
