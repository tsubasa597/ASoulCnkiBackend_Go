package vo

type Status struct {
	Enable  bool  `json:"listen"`
	Started bool  `json:"runing"`
	Wait    int32 `json:"num"`
}
