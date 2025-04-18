use reqwest::Client;
use std::error::Error;

use quick_xml::events::Event;
use quick_xml::reader::Reader;

use tokio::fs::File;
use tokio::io::AsyncWriteExt; // for write_all()

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let url = "https://awesomemotive.com/feed";
    
    let client = Client::new();
    
    let response = client.get(url).send().await?;
    
    if response.status().is_success() {
        let body = response.text().await?; 
        save_feed(&body).await?;    
        println!("File feed.xml saved successfully");
        find(&body);
    } else {
        eprintln!("Error in request: \n {} \n", response.status());
    }
    // parser("feed.xml".to_string())?;
    Ok(())
}
async fn save_feed(text: &String) -> Result<(),Box<dyn Error>>{
    let mut file = File::create("feed.xml").await?;
    file.write_all(text.as_bytes()).await?;

    Ok(())
}


fn find(text: &String){
    for line in text.lines(){
        if line.contains("generator"){
            println!("Version found at line {}", line);
        }
    }
}
