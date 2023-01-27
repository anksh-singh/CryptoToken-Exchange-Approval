package unmarshal

type UnmarshallAssetModel struct {
	ContractName         string  `json:"contract_name"`
	ContractTickerSymbol string  `json:"contract_ticker_symbol"`
	ContractDecimals     int32   `json:"contract_decimals"`
	ContractAddress      string  `json:"contract_address"`
	Coin                 int64   `json:"coin"`
	Type                 string  `json:"type"`
	Balance              string  `json:"balance"`
	Quote                float64 `json:"quote"`
	QuoteRate            float64 `json:"quote_rate"`
	LogoURL              string  `json:"logo_url"`
	QuoteRate24H         string  `json:"quote_rate_24h"`
	QuotePctChange24H    float64 `json:"quote_pct_change_24h"`
}

type UnmarshallTransactionModel struct {
	Page         int64 `json:"page"`
	TotalPages   int64 `json:"total_pages"`
	ItemsOnPage  int64 `json:"items_on_page"`
	TotalTxs     int64 `json:"total_txs"`
	Transactions []struct {
		ID                  string `json:"id"`
		From                string `json:"from"`
		To                  string `json:"to"`
		Fee                 string `json:"fee"`
		Date                int64  `json:"date"`
		Status              string `json:"status"`
		Type                string `json:"type"`
		Block               int64  `json:"block"`
		Value               string `json:"value"`
		Nonce               int64  `json:"nonce"`
		NativeTokenDecimals int64  `json:"native_token_decimals"`
		Description         string `json:"description"`
		Received            []struct {
			Name      string  `json:"name"`
			Symbol    string  `json:"symbol"`
			TokenID   string  `json:"token_id"`
			Decimals  int64   `json:"decimals"`
			Value     string  `json:"value"`
			Quote     float32 `json:"quote"`
			QuoteRate float32 `json:"quoteRate"`
			LogoURL   string  `json:"logo_url"`
			From      string  `json:"from"`
			To        string  `json:"to"`
		} `json:"received,omitempty"`
		Sent []struct {
			Name      string  `json:"name"`
			Symbol    string  `json:"symbol"`
			TokenID   string  `json:"token_id"`
			Decimals  int64   `json:"decimals"`
			Value     string  `json:"value"`
			Quote     float32 `json:"quote"`
			QuoteRate float32 `json:"quoteRate"`
			LogoURL   string  `json:"logo_url"`
			From      string  `json:"from"`
			To        string  `json:"to"`
		} `json:"sent,omitempty"`
		Others []struct {
			Name      string  `json:"name"`
			Symbol    string  `json:"symbol"`
			TokenID   string  `json:"token_id"`
			Decimals  int64   `json:"decimals"`
			Value     string  `json:"value"`
			Quote     float32 `json:"quote"`
			QuoteRate float32 `json:"quoteRate"`
			LogoURL   string  `json:"logo_url"`
			From      string  `json:"from"`
			To        string  `json:"to"`
		} `json:"others,omitempty"`
	} `json:"transactions"`
}

type UserDataModel struct {
	AverageTokenPrice      float64 `json:"average_token_price"`
	CurrentHoldingQuantity float64 `json:"current_holding_quantity"`
	OverallProfitLoss      float64 `json:"overall_profit_loss"`
	PercentageChange24H    float64 `json:"percentage_change_24H"`
	PriceChange24H         float64 `json:"price_change_24H"`
	QuoteRate              float64 `json:"quote_rate"`
	TotalFeesPaid          float64 `json:"total_fees_paid"`
	TotalFeesPaidUsd       float64 `json:"total_fees_paid_usd"`
}

type NFTCollectionDataModel struct {
	Assets []struct {
		Collection struct {
			BannerImageUrl          string `json:"banner_image_url"`
			ChatUrl                 string `json:"chat_url"`
			CreatedDate             string `json:"created_date"`
			DefaultToFiat           bool   `json:"default_to_fiat"`
			Description             string `json:"description"`
			DevBuyerFeeBasisPoints  string `json:"dev_buyer_fee_basis_points"`
			DevSellerFeeBasisPoints string `json:"dev_seller_fee_basis_points"`
			DiscordUrl              string `json:"discord_url"`
			DisplayData             struct {
				CardDisplayStyle string `json:"card_display_style"`
			} `json:"display_data"`
			ExternalUrl                 string `json:"external_url"`
			Featured                    bool   `json:"featured"`
			FeaturedImageUrl            string `json:"featured_image_url"`
			Hidden                      bool   `json:"hidden"`
			SafelistRequestStatus       string `json:"safelist_request_status"`
			ImageUrl                    string `json:"image_url"`
			IsSubjectToWhitelist        bool   `json:"is_subject_to_whitelist"`
			LargeImageUrl               string `json:"large_image_url"`
			MediumUsername              string `json:"medium_username"`
			Name                        string `json:"name"`
			OnlyProxiedTransfers        bool   `json:"only_proxied_transfers"`
			OpenseaBuyerFeeBasisPoints  string `json:"opensea_buyer_fee_basis_points"`
			OpenseaSellerFeeBasisPoints string `json:"opensea_seller_fee_basis_points"`
			PayoutAddress               string `json:"payout_address"`
			RequireEmail                bool   `json:"require_email"`
			ShortDescription            string `json:"short_description"`
			Slug                        string `json:"slug"`
			TelegramUrl                 string `json:"telegram_url"`
			TwitterUsername             string `json:"twitter_username"`
			InstagramUsername           string `json:"instagram_username"`
			WikiUrl                     string `json:"wiki_url"`
			IsNsfw                      bool   `json:"is_nsfw"`
		} `json:"collection"`
		Id                   int64  `json:"id"`
		NumSales             int64  `json:"num_sales"`
		BackgroundColor      string `json:"background_color"`
		ImageUrl             string `json:"image_url"`
		ImagePreviewUrl      string `json:"image_preview_url"`
		ImageThumbnailUrl    string `json:"image_thumbnail_url"`
		ImageOriginalUrl     string `json:"image_original_url"`
		AnimationUrl         string `json:"animation_url"`
		AnimationOriginalUrl string `json:"animation_original_url"`
		Name                 string `json:"name"`
		Description          string `json:"description"`
		ExternalLink         string `json:"external_link"`
		AssetContract        struct {
			Address string `json:"address"`
		} `json:"asset_contract"`
		Permalink     string `json:"permalink"`
		Decimals      string `json:"decimals"`
		TokenMetadata string `json:"token_metadata"`
		IsNsfw        bool   `json:"is_nsfw"`
		Owner         struct {
		} `json:"owner"`
		SellOrders struct {
		} `json:"owner"`
		SeaportSellOrders string `json:"seaport_sell_orders"`
		Creator           struct {
			User struct {
				Username string `json:"username"`
			} `json:"user"`
			ProfileImgUrl string `json:"profile_img_url"`
			Address       string `json:"address"`
			Config        string `json:"config"`
		} `json:"creator"`
		Traits []struct {
			TraitType   string `json:"trait_type"`
			Value       string `json:"value"`
			DisplayType string `json:"display_type"`
			MaxValue    string `json:"max_value"`
			TraitCount  int64  `json:"trait_count"`
			Order       string `json:"order"`
		} `json:"traits"`

		LastSale struct {
			Amount string `json:"order"`
			Symbol string `json:"amount"`
		} `json:"last_sale"`
		TopBid                  string `json:"top_bid"`
		ListingDate             string `json:"listing_date"`
		IsPresale               bool   `json:"is_presale"`
		TransferFeePaymentToken string `json:"transfer_fee_payment_token"`
		TransferFee             string `json:"transfer_fee"`
		TokenId                 string `json:"token_id"`
		CollectionName          string `json:"name"`
		ContractAddress         string `json:"num_sales"`
	} `json:"assets,omitempty"`
}

type EmptyTransactionResponse struct {
	Transactions []interface{} `json:"transactions"`
}
