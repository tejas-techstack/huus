package cache

import (
  types "github.com/tejas-techstack/huus/internal/engine/types"
  errs "github.com/tejas-techstack/huus/internal/handlers/error_handlers"
  page "github.com/tejas-techstack/huus/internal/engine/pager"
)

type Frame struct {
  headers FrameHeaders,
  page    page.Page
}
