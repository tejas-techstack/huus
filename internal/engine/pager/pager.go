package pager

import (
  types "github.com/tejas-techstack/huus/internal/engine/types"
  errs "github.com/tejas-techstack/huus/internal/handlers/error_handlers"
)

/*
  We follow slotted paging technique
  refer : https://siemens.blog/posts/database-page-layout/
*/

type PageType int

// define enum for pageType
const (
  ROOT pageType = iota
  INTERIOR
  LEAF
)

type Page struct{
  pageID uint32
  page PageType
  freeSpace uint32
  fsStart uint32
  fsEnd uint32
}

func NewPage() (Page, error) {}

func (p *Page) AddCell(cell CellPointer) (CellIndex, error) {}

func (p *Page) DeleteCell(idx CellIndex) error {}

func (p *Page) LoadPage() error {}
