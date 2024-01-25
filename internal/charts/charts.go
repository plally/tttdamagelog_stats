package charts

type PieChartItem struct {
	Label string
	Data  int
}

type PieChartInput struct {
	Label string
	Items []PieChartItem
}

func PieChart(input PieChartInput) Chart {
	labels := []string{}
	data := []int{}

	for _, item := range input.Items {
		labels = append(labels, item.Label)
		data = append(data, item.Data)
	}

	return Chart{
		Type:    "pie",
		Options: DefaultOptions(),
		Data: Data{
			Annotation: input.Label,
			Labels:     labels,
			Datasets: []Dataset{
				{
					Label: input.Label,
					Data:  data,
				},
			},
		},
	}
}

type BarChartItem struct {
	Label string
	Data  int
}

type BarChartInput struct {
	Label string
	Items []BarChartItem
}

func BarChart(input BarChartInput) Chart {
	labels := []string{}
	data := []int{}

	for _, item := range input.Items {
		labels = append(labels, item.Label)
		data = append(data, item.Data)
	}

	return Chart{
		Type:    "bar",
		Options: DefaultOptions(),
		Data: Data{
			Annotation: input.Label,
			Labels:     labels,
			Datasets: []Dataset{
				{
					Label: input.Label,
					Data:  data,
				},
			},
		},
	}
}

func DefaultOptions() Options {
	return Options{
		Events:              []string{},
		MaintainAspectRatio: false,
	}
}

type Chart struct {
	Type    string  `json:"type"`
	Data    Data    `json:"data"`
	Options Options `json:"options"`
}

type Data struct {
	Annotation string    `json:"annotation"`
	Labels     []string  `json:"labels"`
	Datasets   []Dataset `json:"datasets"`
}

type Dataset struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

type Options struct {
	Events              []string `json:"events"`
	MaintainAspectRatio bool     `json:"maintainAspectRatio"`
}
