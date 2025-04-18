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
        match parser(&body){
            Ok(content) => println!("Wordpress generator content:\n {}", content),
            Err(e) => {println!("Error while getting content, trying in another way... \n Error:\t{}", e);
                    find(&body);
            }
        }
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

// if parser fails
fn find(text: &String){
    for line in text.lines(){
        if line.contains("generator"){
            println!("Version found at line {}", line);
        }
    }
}


fn parser(xml: &str) -> Result<String,Box<dyn Error>>{
    let mut reader = Reader::from_str(xml);
    // reader.trim_text(true);

    let mut buf = Vec::new();
    let mut inside_tag = false;

    loop{
        match reader.read_event_into(&mut buf){
            Ok(Event::Start(ref e)) if e.name().as_ref() == b"generator" => {
                inside_tag = true;
            }

            Ok(Event::Text(e)) if inside_tag => {
                let text = e.unescape().unwrap_or_default().to_string();
                return Ok(text);
            }

            Ok(Event::Eof) => break,

            Ok(Event::End(ref e)) if e.name().as_ref() == b"generator" => {
                inside_tag = false;
            }

            Err(e) => {
                return Err(Box::new(e));
            }
            
            _ => {}
        }

        buf.clear();
    }
    
    Err("No <generator> tag found".into())
}
