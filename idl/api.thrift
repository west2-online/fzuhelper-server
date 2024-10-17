namespace go api
include "model.thrift"

//user
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required string id
    2: required list<string> cookies
}

struct EmptyClassroomRequest {
    1: required string date
    2: required string campus
    3: required string startTime;//节数
    4: required string endTime;
}

struct EmptyClassroomResponse {
    1: required list<model.Classroom> classrooms
}



//paper
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
    1: required model.UpYunFileDir dir,
}

struct GetDownloadUrlRequest {
    1: required string url,
}

struct GetDownloadUrlResponse {
    1: required string url,
}


service ClassRoomService {
    EmptyClassroomResponse GetEmptyClassrooms(1: EmptyClassroomRequest request)(api.get="/api/v1/common/classroom/empty")
}

service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest request)(api.get="/api/v1/jwch/user/login")
}

service PaperService {
    UploadFileResponse UploadFile(1: UploadFileRequest req) (api.post="/api/v1/paper/upload"),
    ListDirFilesResponse ListDirFiles(1: ListDirFilesRequest req) (api.get="/api/v1/paper/list"),
    GetDownloadUrlResponse GetDownloadUrl(1: GetDownloadUrlRequest req) (api.get="/api/v1/paper/download"),
}

