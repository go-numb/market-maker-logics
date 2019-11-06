// 売り -- LTP -- 買い
// 1038361 -- 1038163 -- 1037821,  0.999930(last: 1.000004, past: 1.000039)
// Ind: 1.000035
// Go trade, buy & sell orders
// 売り -- LTP -- 買い
// 1038431 -- 1038172 -- 1037891,  0.999989(last: 0.999993, past: 0.999998)
// Ind: 1.000005
// ---------------------------


// Reset is create mean
func (p *Execute) Reset(spanVolume float64) { // 乖離加速度Logic@taro
	p.Lock()
	defer p.Unlock()

	devAvg := p.Price / stat.Mean(p.Prices, nil) // 直近n秒乖離

	// 現在を過去n秒の配列に追加
	l := len(p.PricesPast)
	p.PricesPast = append(p.PricesPast, p.Prices...) // 過去に現在を追加
	p.SizesPast = append(p.SizesPast, p.Sizes...)    // 過去に現在を追加
	if !math.IsNaN(devAvg) {                         // !IsZero ならば乖離平均配列へ追加
		p.DevAvgs = append(p.DevAvgs, devAvg)
	}

	// 直近と過去配列から乖離平均を算出
	devAvgPast := stat.Mean(p.DevAvgs, nil) // n秒「乖離」の平均

	// 直近の最高値・最安値から想定値動きを取得して0.0026を変えるムーブは状態が変動するので改善点が追えなくなる
	// よって、過去を見て変数を変えるのは排除

	if 0 < l { // 不要部分を削除
		p.PricesPast = p.PricesPast[l:]
		p.SizesPast = p.SizesPast[l:]
		p.DevAvgs = p.DevAvgs[len(p.DevAvgs)-1:]
		p.Prices = []float64{}
		p.Sizes = []float64{}
	}

	sma := ((devAvg - devAvgPast) + devAvg) / devAvgPast

	p.AskDistance = p.Price * (sma + RANGEPARPRICE)
	p.BidDistance = p.Price * (sma - RANGEPARPRICE)

	// (乖離平均 > 1(価格上昇傾向) and 指標 > 1(ltp上昇してない))
	// or
	// (乖離平均 < 1(価格下落傾向) and 指標 < 1(ltp下落してない))
	ind := devAvgPast / devAvg
	// fmt.Println("売り -- LTP -- 買い")
	// fmt.Printf("%.f -- %.f -- %.f,	%f(last: %f, past: %f)\n", p.AskDistance, p.Price, p.BidDistance, sma, devAvg, devAvgPast)
	// fmt.Printf("Ind: %f\n", ind)

	if 1 < devAvgPast && 1 < ind ||
		devAvgPast < 1 && ind < 1 {
		fmt.Println("Go trade, buy & sell orders")
	} else {
		fmt.Println("---------------------------")
	}

	p.Volume = 0
}