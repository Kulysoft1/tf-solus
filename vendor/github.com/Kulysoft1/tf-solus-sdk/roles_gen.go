// Autogenerated file. Do not edit!

package solus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Next using for iterating through all data entities.
//
// Examples:
//
//	ctx, cancelFunc := context.WithTimeout(context.Background(), 30 * time.Second)
//	defer cancelFunc()
//
//  for {
//		for _, datum := range resp.Data {
//			doSmthWithDatum(datum)
//		}
//
//		if !resp.Next(ctx) {
//			break
//		}
//	}
//
//  if resp.Err() != nil {
//		handleAnError(resp.Err())
//	}
func (r *RolesResponse) Next(ctx context.Context) bool {
	if (r.Meta.LastPage == r.Meta.CurrentPage) || (r.err != nil) {
		return false
	}

	body, code, err := r.service.client.request(ctx, http.MethodGet, r.Links.Next)
	if err != nil {
		r.err = err
		return false
	}

	if code != http.StatusOK {
		r.err = newHTTPError(http.MethodGet, r.Links.Next, code, body)
		return false
	}

	if err := json.Unmarshal(body, &r); err != nil {
		r.err = fmt.Errorf("failed to decode %q: %s", body, err)
		return false
	}
	return true
}