// Copyright (c) 2016 turbobytes
//
// This file is part of geoipdb, a library of GeoIP related helper functions
// for TurboBytes stack.
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package geoipdb

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AsnOverride is what is stored in the overrides collection.
type AsnOverride struct {
	Asn  string `bson:"_id" json:"asn"`
	Name string `bson:"name" json:"name"`
}

// OverridesNilCollectionError is returned by Overrides<...> methods
// when Handler was created without an overrides collection
// (see NewHandler).
var OverridesNilCollectionError = errors.New("nil overrides collection")

// OverridesAsnNotFoundError is returned by OverridesLookup
// when there is no override defined.
var OverridesAsnNotFoundError = errors.New("ASN not found")

// OverridesMalformedAsnError is returned by OverridesSet
// when parameter asn does not conform to an ASN identification.
var OverridesMalformedAsnError = errors.New("malformed ASN")

// OverridesLookup queries the database of local overrides
// for the description of a given ASN.
//
// Returns the ASN description,
// or OverridesAsnNotFoundError if there is no override for the ASN.
func (h Handler) OverridesLookup(asn string) (string, error) {
	if h.overrides == nil {
		return "", OverridesNilCollectionError
	}
	var override AsnOverride
	err := h.overrides.FindId(asn).One(&override)
	if err == mgo.ErrNotFound {
		return "", OverridesAsnNotFoundError
	}
	if err != nil {
		return "", fmt.Errorf("cannot lookup override: %s", err)
	}
	return override.Name, nil
}

// OverridesSet stores or updates a user defined description for a given ASN
// in the database of local overrides.
//
// Moreover, this method purges the cache (see LookupAsn)
// of all data related to the given asn.
func (h Handler) OverridesSet(asn string, descr string) error {
	h.cache.purgeASN(asn)
	if h.overrides == nil {
		return OverridesNilCollectionError
	}
	if !reASN.MatchString(asn) {
		return OverridesMalformedAsnError
	}
	_, err := h.overrides.UpsertId(asn, bson.M{"$set": bson.M{"name": descr}})
	if err != nil {
		return fmt.Errorf("cannot set override: %s", err)
	}
	return nil
}

// OverridesRemove removes the description for a given ASN
// from the database of local overrides.
// If there is no such ASN,
// OverridesRemove returns silently without error.
//
// Moreover, this method purges the cache (see LookupAsn)
// of all data related to the given asn.
func (h Handler) OverridesRemove(asn string) error {
	h.cache.purgeASN(asn)
	if h.overrides == nil {
		return OverridesNilCollectionError
	}
	err := h.overrides.RemoveId(asn)
	if err == mgo.ErrNotFound {
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot remove override: %s", err)
	}
	return nil
}

// OverridesList answers all ASN description overrides.
func (h Handler) OverridesList() ([]AsnOverride, error) {
	if h.overrides == nil {
		return nil, OverridesNilCollectionError
	}
	var answer []AsnOverride
	err := h.overrides.Find(nil).All(&answer)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve overrides: %s", err)
	}
	if answer == nil {
		return make([]AsnOverride, 0), nil
	}
	return answer, nil
}
