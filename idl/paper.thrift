namespace go paper

include "model.thrift"




struct UploadFileRequest {
    1: required string filename,
    2: required binary content,
    3: required string ussPath,
    4: required string user,
}

struct UploadFileResponse {
    1: required model.BaseResp base,
}

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
    UploadFileResponse UploadFile(1: UploadFileRequest req),
    ListDirFilesResponse ListDirFiles(1: ListDirFilesRequest req),
    GetDownloadUrlResponse GetDownloadUrl(1: GetDownloadUrlRequest req),
}

