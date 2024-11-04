namespace go paper

include "model.thrift"

struct ListDirFilesRequest {
    1: required string path,

}

struct ListDirFilesResponse {
    1: required model.BaseResp base,
    2: required model.UpYunFileDir dir,
}

struct GetDownloadUrlRequest {
    1: required string url,
}

struct GetDownloadUrlResponse {
    1: required model.BaseResp base,
    2: required string url,
}

service PaperService {
    ListDirFilesResponse ListDirFiles(1: ListDirFilesRequest req),
    GetDownloadUrlResponse GetDownloadUrl(1: GetDownloadUrlRequest req),
}

