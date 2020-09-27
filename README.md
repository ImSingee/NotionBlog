# Notion + Blog = NB!
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FImSingee%2FNotionBlog.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FImSingee%2FNotionBlog?ref=badge_shield)


Use Notion to write, use hexo to deploy.

[中文文档](https://www.notion.so/singee/NotionBlog-44f5de5864fa4ef19dda4d7f57ab3652)

## Usage

You need to do three things:

1. Create a notion database to put your posts
2. Make a `_notion` folder in your hexo's `source` folder, and create a `config.yml` file
3. Download this project's file and run `./nb -root "your hexo's root path"`


## Config file's format

You need a `config.yml` file in your hexo's `source/_notion` folder. The config file must contain the following parts

```yaml
version: 1 # Now should be fixed to 1
token_v2: # Your notion's token
database:
  post:
  - 963f630adc2e443b98c7c93378c17176+4312fa9b8f8142a0832a95008cfee6c0

```

And there're also some optional settings
```yaml
converter:
  force: default # set to true to rerender all pages, otherwise only rerender edited files
render:
  checkbox: false # set to true to render "To-do" block to checkbox, otherwise to normal list
user:
  locale: en
  timezone: Etc/UTC # tz database time zones
```




## License

This software is released under the Apache-2.0 license.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FImSingee%2FNotionBlog.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FImSingee%2FNotionBlog?ref=badge_large)