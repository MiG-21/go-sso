package types

type HealthCheckPing struct {
	Ping string `json:"ping"`
}

type HealthCheckInfoGit struct {
	Hash string `json:"hash"`
	Ref  string `json:"ref"`
	Url  string `json:"url"`
}

type HealthCheckInfo struct {
	AppName     string             `json:"appName"`
	AppVersion  string             `json:"appVersion"`
	ClusterName string             `json:"clusterName"`
	Git         HealthCheckInfoGit `json:"git"`
}
