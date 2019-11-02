package taros

import "fmt"

// Reset is  resets by ticker time.
// Use executions ticker or executions, and volumes to vwap.
func (p *Execute) Reset() { // 乖離加速度Logic@taro
	mean := stat.Mean(p.Prices, p.Sizes) // 直近n秒SMA
	l := len(p.PricesPast)
	p.PricesPast = append(p.PricesPast, p.Prices...) // 過去に現在を追加
	p.SizesPast = append(p.SizesPast, p.Sizes...)    // 過去に現在を追加
	meanLast := stat.Mean(p.PricesPast, p.SizesPast) // n+x足SMA
	if 0 < l {                                       // 不要部分を削除
		p.PricesPast = p.PricesPast[l:]
		p.SizesPast = p.SizesPast[l:]
		p.Prices = []float64{}
		p.Sizes = []float64{}
	}

	sma := ((mean - meanLast) + mean) / meanLast

	askPrice := mean * (sma + 0.00026)
	bidPrice := mean * (sma - 0.00026)
	diffBuy := askPrice - mean
	diffSell := mean - bidPrice
	if diffBuy < diffSell {
		p.IsBuyByAcceleration = false
	} else if diffBuy > diffSell {
		p.IsBuyByAcceleration = true
	}

	p.PriceDivergense = mean
	// fmt.Printf("%t\n", p.IsBuyByAcceleration)
	// fmt.Printf("---------------- %.f (%.f/%.f/%f) %.f ++ %.f  (%.f)\n", p.Price, mean, meanLast, sma, diffBuy, diffSell, bidPrice-askPrice)
}
