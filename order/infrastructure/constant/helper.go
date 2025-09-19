package constant

const (
	OrderStatusCreated    = 0
	OrderStatusProcessing = 1
	OrderStatusCompleted  = 2
	OrderStatusCancelled  = 3
)

var OrderStatusTranslated = map[int]string{
	OrderStatusCreated:    "Created",
	OrderStatusProcessing: "Processing",
	OrderStatusCompleted:  "Completed",
	OrderStatusCancelled:  "Cancelled",
}
