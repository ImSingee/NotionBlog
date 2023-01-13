use notion::ids::{AsIdentifier, BlockId, PageId};
use notion::{Error, NotionApi};

#[derive(Debug, Eq, PartialEq, Clone)]
struct Page {
    pub meta: notion::models::Page,
}

#[async_trait::async_trait]
trait GetFullPage {
    async fn get_full_page<T: AsIdentifier<PageId> + std::marker::Send>(&self, page_id: T) -> Result<Page, Error>;
}

#[async_trait::async_trait]
impl GetFullPage for NotionApi {
    async fn get_full_page<T: AsIdentifier<PageId> + std::marker::Send>(&self, page_id: T) -> Result<Page, Error> {
        let meta = self.get_page(page_id.as_id()).await?;

        let block_id = BlockId::from(page_id);


        let blocks = self.get_block_children(block_id).await?;
        for block in blocks.results {

        }


        Ok(Page {
            meta
        })
    }
}