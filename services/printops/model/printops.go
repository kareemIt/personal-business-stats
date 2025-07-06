package model

type PrintJob struct {
	PriceCost      int
	TrackingNumber string
}

func NewPrintJob(priceCost int, trackingNumber string) *PrintJob {
	if priceCost == 0 || trackingNumber == "" {
		panic("Price cost cannot be zero and tracking number cannot be empty")
	}
	return &PrintJob{
		PriceCost:      priceCost,
		TrackingNumber: trackingNumber,
	}
}

func (p *PrintJob) GetPriceCost() int {
	return p.PriceCost
}

func (p *PrintJob) GetTrackingNumber() string {
	return p.TrackingNumber
}

func (p *PrintJob) SetPriceCost(priceCost int) {
	if priceCost == 0 {
		panic("Price cost cannot be zero")
	}
	p.PriceCost = priceCost
}
func (p *PrintJob) SetTrackingNumber(trackingNumber string) {
	if trackingNumber == "" {
		panic("Tracking number cannot be empty")
	}
	p.TrackingNumber = trackingNumber
}
