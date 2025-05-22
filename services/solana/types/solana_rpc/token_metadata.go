package solana_types

type GetTokenMetaDataResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  TokenMetaData `json:"result"`
	ID      int           `json:"id"`
}

type TokenMetaData struct {
	Interface   string        `json:"interface"`
	ID          string        `json:"id"`
	Content     Content       `json:"content"`
	Authorities []Authority   `json:"authorities"`
	Compression Compression   `json:"compression"`
	Grouping    []interface{} `json:"grouping"`
	Royalty     Royalty       `json:"royalty"`
	Creators    []interface{} `json:"creators"`
	Ownership   Ownership     `json:"ownership"`
	Mutable     bool          `json:"mutable"`
	Burnt       bool          `json:"burnt"`
}

type Content struct {
	Schema   string   `json:"$schema"`
	JSONURI  string   `json:"json_uri"`
	Files    []File   `json:"files"`
	Metadata Metadata `json:"metadata"`
	Links    Links    `json:"links"`
}

type File struct {
	URI  string `json:"uri"`
	Mime string `json:"mime"`
}

type Metadata struct {
	Description   string `json:"description"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	TokenStandard string `json:"token_standard"`
}

type Links struct {
	Image string `json:"image"`
}

type Authority struct {
	Address string   `json:"address"`
	Scopes  []string `json:"scopes"`
}

type Compression struct {
	Eligible    bool   `json:"eligible"`
	Compressed  bool   `json:"compressed"`
	DataHash    string `json:"data_hash"`
	CreatorHash string `json:"creator_hash"`
	AssetHash   string `json:"asset_hash"`
	Tree        string `json:"tree"`
	Seq         int    `json:"seq"`
	LeafID      int    `json:"leaf_id"`
}

type Royalty struct {
	RoyaltyModel        string  `json:"royalty_model"`
	Target              *string `json:"target"`
	Percent             float64 `json:"percent"`
	BasisPoints         int     `json:"basis_points"`
	PrimarySaleHappened bool    `json:"primary_sale_happened"`
	Locked              bool    `json:"locked"`
}

type Ownership struct {
	Frozen         bool    `json:"frozen"`
	Delegated      bool    `json:"delegated"`
	Delegate       *string `json:"delegate"`
	OwnershipModel string  `json:"ownership_model"`
	Owner          string  `json:"owner"`
}
