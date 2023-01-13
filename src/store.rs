use std::collections::{HashMap, HashSet};
use queues::{IsQueue, Queue};
use anyhow::*;
use notion::ids::PageId;
use notion::models::Page;
use notion::NotionApi;

type Map = HashMap<PageId, Page>;

struct Store<'a> {
    notion: &'a NotionApi,
    map: Map,
}

impl<'a> Store<'a> {
    pub fn get(&mut self, id: &PageId) -> Option<&Page> {
        self.map.get(id)
    }

    pub async fn get_or_fetch(&mut self, id: &PageId) -> Result<&Page> {
        if !self.map.contains_key(id) {
            let page = self.notion.get_page(id).await?;
            self.map.insert(id.to_owned(), page);
        }


        self.get(id).ok_or(anyhow!("unexpected error: cannot found page"))
    }

    pub async fn fetch_all(&mut self, ids: &[&PageId]) -> Result<()> {
        let mut queue = Queue::default();
        let mut set = HashSet::with_capacity(ids.len()*2);

        for &id in ids {
            queue.add(id);
            set.insert(id);
        }

        while let std::result::Result::Ok(id) = queue.remove() {
            if !set.contains(id) {
                let page = self.get_or_fetch(id).await?;
            }
        }

        todo!()
    }
}