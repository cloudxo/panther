// Code generated by go-swagger; DO NOT EDIT.

package operations

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
)

// ModifyGlobalReader is a Reader for the ModifyGlobal structure.
type ModifyGlobalReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ModifyGlobalReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewModifyGlobalOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewModifyGlobalBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewModifyGlobalNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewModifyGlobalInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewModifyGlobalOK creates a ModifyGlobalOK with default headers values
func NewModifyGlobalOK() *ModifyGlobalOK {
	return &ModifyGlobalOK{}
}

/*ModifyGlobalOK handles this case with default header values.

OK
*/
type ModifyGlobalOK struct {
	Payload *models.Global
}

func (o *ModifyGlobalOK) Error() string {
	return fmt.Sprintf("[POST /global/update][%d] modifyGlobalOK  %+v", 200, o.Payload)
}

func (o *ModifyGlobalOK) GetPayload() *models.Global {
	return o.Payload
}

func (o *ModifyGlobalOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Global)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewModifyGlobalBadRequest creates a ModifyGlobalBadRequest with default headers values
func NewModifyGlobalBadRequest() *ModifyGlobalBadRequest {
	return &ModifyGlobalBadRequest{}
}

/*ModifyGlobalBadRequest handles this case with default header values.

Bad request
*/
type ModifyGlobalBadRequest struct {
	Payload *models.Error
}

func (o *ModifyGlobalBadRequest) Error() string {
	return fmt.Sprintf("[POST /global/update][%d] modifyGlobalBadRequest  %+v", 400, o.Payload)
}

func (o *ModifyGlobalBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *ModifyGlobalBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewModifyGlobalNotFound creates a ModifyGlobalNotFound with default headers values
func NewModifyGlobalNotFound() *ModifyGlobalNotFound {
	return &ModifyGlobalNotFound{}
}

/*ModifyGlobalNotFound handles this case with default header values.

Global not found
*/
type ModifyGlobalNotFound struct {
}

func (o *ModifyGlobalNotFound) Error() string {
	return fmt.Sprintf("[POST /global/update][%d] modifyGlobalNotFound ", 404)
}

func (o *ModifyGlobalNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewModifyGlobalInternalServerError creates a ModifyGlobalInternalServerError with default headers values
func NewModifyGlobalInternalServerError() *ModifyGlobalInternalServerError {
	return &ModifyGlobalInternalServerError{}
}

/*ModifyGlobalInternalServerError handles this case with default header values.

Internal server error
*/
type ModifyGlobalInternalServerError struct {
}

func (o *ModifyGlobalInternalServerError) Error() string {
	return fmt.Sprintf("[POST /global/update][%d] modifyGlobalInternalServerError ", 500)
}

func (o *ModifyGlobalInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
