/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by thriftgo (0.3.17). DO NOT EDIT.

package model

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type BaseResp struct {
	Code int64  `thrift:"code,1" form:"code" json:"code" query:"code"`
	Msg  string `thrift:"msg,2" form:"msg" json:"msg" query:"msg"`
}

func NewBaseResp() *BaseResp {
	return &BaseResp{}
}

func (p *BaseResp) InitDefault() {
}

func (p *BaseResp) GetCode() (v int64) {
	return p.Code
}

func (p *BaseResp) GetMsg() (v string) {
	return p.Msg
}

var fieldIDToName_BaseResp = map[int16]string{
	1: "code",
	2: "msg",
}

func (p *BaseResp) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.I64 {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 2:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}
		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_BaseResp[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *BaseResp) ReadField1(iprot thrift.TProtocol) error {

	var _field int64
	if v, err := iprot.ReadI64(); err != nil {
		return err
	} else {
		_field = v
	}
	p.Code = _field
	return nil
}
func (p *BaseResp) ReadField2(iprot thrift.TProtocol) error {

	var _field string
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		_field = v
	}
	p.Msg = _field
	return nil
}

func (p *BaseResp) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("BaseResp"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}
	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *BaseResp) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("code", thrift.I64, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteI64(p.Code); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *BaseResp) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("msg", thrift.STRING, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Msg); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *BaseResp) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaseResp(%+v)", *p)

}

// 由前端给的登陆信息，包括id和cookies, 这个struct仅用于测试返回数据，因为登录实现在前端完成，不会在实际项目中使用
type LoginData struct {
	ID      string   `thrift:"id,1,required" form:"id,required" json:"id,required" query:"id,required"`
	Cookies []string `thrift:"cookies,2,required" form:"cookies,required" json:"cookies,required" query:"cookies,required"`
}

func NewLoginData() *LoginData {
	return &LoginData{}
}

func (p *LoginData) InitDefault() {
}

func (p *LoginData) GetID() (v string) {
	return p.ID
}

func (p *LoginData) GetCookies() (v []string) {
	return p.Cookies
}

var fieldIDToName_LoginData = map[int16]string{
	1: "id",
	2: "cookies",
}

func (p *LoginData) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16
	var issetID bool = false
	var issetCookies bool = false

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
				issetID = true
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 2:
			if fieldTypeId == thrift.LIST {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
				issetCookies = true
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}
		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	if !issetID {
		fieldId = 1
		goto RequiredFieldNotSetError
	}

	if !issetCookies {
		fieldId = 2
		goto RequiredFieldNotSetError
	}
	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_LoginData[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
RequiredFieldNotSetError:
	return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("required field %s is not set", fieldIDToName_LoginData[fieldId]))
}

func (p *LoginData) ReadField1(iprot thrift.TProtocol) error {

	var _field string
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		_field = v
	}
	p.ID = _field
	return nil
}
func (p *LoginData) ReadField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return err
	}
	_field := make([]string, 0, size)
	for i := 0; i < size; i++ {

		var _elem string
		if v, err := iprot.ReadString(); err != nil {
			return err
		} else {
			_elem = v
		}

		_field = append(_field, _elem)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return err
	}
	p.Cookies = _field
	return nil
}

func (p *LoginData) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("LoginData"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}
	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *LoginData) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("id", thrift.STRING, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.ID); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *LoginData) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("cookies", thrift.LIST, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteListBegin(thrift.STRING, len(p.Cookies)); err != nil {
		return err
	}
	for _, v := range p.Cookies {
		if err := oprot.WriteString(v); err != nil {
			return err
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *LoginData) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LoginData(%+v)", *p)

}

// Classroom 前端想要返回的fields
type Classroom struct {
	Build    string `thrift:"build,1,required" form:"build,required" json:"build,required" query:"build,required"`
	Location string `thrift:"location,2,required" form:"location,required" json:"location,required" query:"location,required"`
	Capacity string `thrift:"capacity,3,required" form:"capacity,required" json:"capacity,required" query:"capacity,required"`
	Type     string `thrift:"type,4,required" form:"type,required" json:"type,required" query:"type,required"`
}

func NewClassroom() *Classroom {
	return &Classroom{}
}

func (p *Classroom) InitDefault() {
}

func (p *Classroom) GetBuild() (v string) {
	return p.Build
}

func (p *Classroom) GetLocation() (v string) {
	return p.Location
}

func (p *Classroom) GetCapacity() (v string) {
	return p.Capacity
}

func (p *Classroom) GetType() (v string) {
	return p.Type
}

var fieldIDToName_Classroom = map[int16]string{
	1: "build",
	2: "location",
	3: "capacity",
	4: "type",
}

func (p *Classroom) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16
	var issetBuild bool = false
	var issetLocation bool = false
	var issetCapacity bool = false
	var issetType bool = false

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
				issetBuild = true
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 2:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
				issetLocation = true
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 3:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField3(iprot); err != nil {
					goto ReadFieldError
				}
				issetCapacity = true
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 4:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField4(iprot); err != nil {
					goto ReadFieldError
				}
				issetType = true
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}
		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	if !issetBuild {
		fieldId = 1
		goto RequiredFieldNotSetError
	}

	if !issetLocation {
		fieldId = 2
		goto RequiredFieldNotSetError
	}

	if !issetCapacity {
		fieldId = 3
		goto RequiredFieldNotSetError
	}

	if !issetType {
		fieldId = 4
		goto RequiredFieldNotSetError
	}
	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_Classroom[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
RequiredFieldNotSetError:
	return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("required field %s is not set", fieldIDToName_Classroom[fieldId]))
}

func (p *Classroom) ReadField1(iprot thrift.TProtocol) error {

	var _field string
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		_field = v
	}
	p.Build = _field
	return nil
}
func (p *Classroom) ReadField2(iprot thrift.TProtocol) error {

	var _field string
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		_field = v
	}
	p.Location = _field
	return nil
}
func (p *Classroom) ReadField3(iprot thrift.TProtocol) error {

	var _field string
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		_field = v
	}
	p.Capacity = _field
	return nil
}
func (p *Classroom) ReadField4(iprot thrift.TProtocol) error {

	var _field string
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		_field = v
	}
	p.Type = _field
	return nil
}

func (p *Classroom) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("Classroom"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}
		if err = p.writeField3(oprot); err != nil {
			fieldId = 3
			goto WriteFieldError
		}
		if err = p.writeField4(oprot); err != nil {
			fieldId = 4
			goto WriteFieldError
		}
	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *Classroom) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("build", thrift.STRING, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Build); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *Classroom) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("location", thrift.STRING, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Location); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *Classroom) writeField3(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("capacity", thrift.STRING, 3); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Capacity); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 end error: ", p), err)
}

func (p *Classroom) writeField4(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("type", thrift.STRING, 4); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Type); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 4 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 4 end error: ", p), err)
}

func (p *Classroom) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Classroom(%+v)", *p)

}

type UpYunFileDir struct {
	BasePath *string  `thrift:"basePath,1,optional" form:"basePath" json:"basePath,omitempty" query:"basePath"`
	Files    []string `thrift:"files,2,required" form:"files,required" json:"files,required" query:"files,required"`
	Folders  []string `thrift:"folders,3,required" form:"folders,required" json:"folders,required" query:"folders,required"`
}

func NewUpYunFileDir() *UpYunFileDir {
	return &UpYunFileDir{}
}

func (p *UpYunFileDir) InitDefault() {
}

var UpYunFileDir_BasePath_DEFAULT string

func (p *UpYunFileDir) GetBasePath() (v string) {
	if !p.IsSetBasePath() {
		return UpYunFileDir_BasePath_DEFAULT
	}
	return *p.BasePath
}

func (p *UpYunFileDir) GetFiles() (v []string) {
	return p.Files
}

func (p *UpYunFileDir) GetFolders() (v []string) {
	return p.Folders
}

var fieldIDToName_UpYunFileDir = map[int16]string{
	1: "basePath",
	2: "files",
	3: "folders",
}

func (p *UpYunFileDir) IsSetBasePath() bool {
	return p.BasePath != nil
}

func (p *UpYunFileDir) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16
	var issetFiles bool = false
	var issetFolders bool = false

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 2:
			if fieldTypeId == thrift.LIST {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
				issetFiles = true
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 3:
			if fieldTypeId == thrift.LIST {
				if err = p.ReadField3(iprot); err != nil {
					goto ReadFieldError
				}
				issetFolders = true
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}
		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	if !issetFiles {
		fieldId = 2
		goto RequiredFieldNotSetError
	}

	if !issetFolders {
		fieldId = 3
		goto RequiredFieldNotSetError
	}
	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_UpYunFileDir[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
RequiredFieldNotSetError:
	return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("required field %s is not set", fieldIDToName_UpYunFileDir[fieldId]))
}

func (p *UpYunFileDir) ReadField1(iprot thrift.TProtocol) error {

	var _field *string
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		_field = &v
	}
	p.BasePath = _field
	return nil
}
func (p *UpYunFileDir) ReadField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return err
	}
	_field := make([]string, 0, size)
	for i := 0; i < size; i++ {

		var _elem string
		if v, err := iprot.ReadString(); err != nil {
			return err
		} else {
			_elem = v
		}

		_field = append(_field, _elem)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return err
	}
	p.Files = _field
	return nil
}
func (p *UpYunFileDir) ReadField3(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return err
	}
	_field := make([]string, 0, size)
	for i := 0; i < size; i++ {

		var _elem string
		if v, err := iprot.ReadString(); err != nil {
			return err
		} else {
			_elem = v
		}

		_field = append(_field, _elem)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return err
	}
	p.Folders = _field
	return nil
}

func (p *UpYunFileDir) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("UpYunFileDir"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}
		if err = p.writeField3(oprot); err != nil {
			fieldId = 3
			goto WriteFieldError
		}
	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *UpYunFileDir) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetBasePath() {
		if err = oprot.WriteFieldBegin("basePath", thrift.STRING, 1); err != nil {
			goto WriteFieldBeginError
		}
		if err := oprot.WriteString(*p.BasePath); err != nil {
			return err
		}
		if err = oprot.WriteFieldEnd(); err != nil {
			goto WriteFieldEndError
		}
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *UpYunFileDir) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("files", thrift.LIST, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteListBegin(thrift.STRING, len(p.Files)); err != nil {
		return err
	}
	for _, v := range p.Files {
		if err := oprot.WriteString(v); err != nil {
			return err
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *UpYunFileDir) writeField3(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("folders", thrift.LIST, 3); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteListBegin(thrift.STRING, len(p.Folders)); err != nil {
		return err
	}
	for _, v := range p.Folders {
		if err := oprot.WriteString(v); err != nil {
			return err
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 end error: ", p), err)
}

func (p *UpYunFileDir) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UpYunFileDir(%+v)", *p)

}
