package structure

type Yml_catalog struct {
	Shop Shop `xml:"shop"`
}

type Shop struct {
	Name                  string        `xml:"name"`
	Company               string        `xml:"company"`
	Url                   string        `xml:"url"`
	Enable_auto_discounts string        `xml:"enable_auto_discounts"`
	Currencies            Currencies    `xml:"currencies"`
	Categories            Categories    `xml:"categories"`
	Delivery_opts         Delivery_opts `xml:"delivery-options"`
	Offers                Offers        `xml:"offers"`
}

type Currencies struct {
	Currencies []Currency `xml:"currency"`
}

type Currency struct {
	Id   string `xml:"id,attr"`
	Rate string `xml:"rate,attr"`
}

type Categories struct {
	Categories []Category `xml:"category"`
}

type Category struct {
	Id       string `xml:"id,attr"`
	ParentId string `xml:"parentId,attr"`
	Name     string `xml:",chardata"`
}

type Delivery_opts struct {
	Options []Option `xml:"option"`
}

type Option struct {
	Cost         string `xml:"cost,attr"`
	Days         string `xml:"days,attr"`
	Order_before string `xml:"order-before,attr"`
}

type Offers struct {
	Offers []Offer `xml:"offer"`
}

type Offer struct {
	Id          string   `xml:"id,attr"`
	Available   string   `xml:"available,attr"`
	Url         string   `xml:"url"`
	CategoryId  string   `xml:"categoryId"`
	Picture     []string `xml:"picture"`
	Name        string   `xml:"name"`
	Vendor      string   `xml:"vendor"`
	VendorCode  string   `xml:"vendorCode"`
	Weight      string   `xml:"weight"`
	Description string   `xml:"description"`
	Barcodes    []string `xml:"barcode"`
	Dimensions  string   `xml:"dimensions"`
	Model       string   `xml:"model"`
	TypePrefix  string   `xml:"typePrefix"`
	Params      []Param  `xml:"param"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Unit  string `xml:"unit,attr"`
	Value string `xml:",chardata"`
}
