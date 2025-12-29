package busi

func (busi *Busi)DealDownloadOSSFile(key string){
	busi.HuaweiOSSClient.GetObjectSDK(key)
}