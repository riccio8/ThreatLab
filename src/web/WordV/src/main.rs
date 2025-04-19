use reqwest::Client;
use std::error::Error;

use regex::Regex;

use quick_xml::events::Event;
use quick_xml::reader::Reader;

use tokio::fs::File;
use tokio::io::AsyncWriteExt; // for write_all()

use webbrowser;

use std::process::Command;
use std::env::consts::OS;


#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let url = "https://awesomemotive.com/feed";
    
    let client = Client::new();
    
    let response = client.get(url).send().await?;
    
    let mut version: String = String::new();

    if response.status().is_success() {
        let body = response.text().await?; 
        save_feed(&body).await?;    
        println!("File feed.xml saved successfully");
        version = match parser(&body){
            Ok(content) => content,
            Err(_) => find(&body).to_string(),
        };
    } else {
        eprintln!("Error in request: \n {} \n", response.status());
    }
    // parser("feed.xml".to_string())?;
    
    // let version1 = version.retain(|s| s != "")
    
    
    let re = Regex::new(r"https?://(?P<name>[^\.]+)\.org/\?v=(?P<version>[\d\.]+)").unwrap();

    let mut version1 = String::new();

    if let Some(caps) = re.captures(&version) {
        let name = &caps["name"];
        let version = &caps["version"];

        version1 = format!("{}-v={}", name, version);
    }
    let url= format!("https://cve.mitre.org/cgi-bin/cvekey.cgi?keyword={}", version1);
        
    
    println!("{}", version1);

    match openBrowser(&url) {
        Ok(_) => {}
        Err(e) => {
            eprintln!("Trying in another way beacuse of error {}...", e);
            if let Err(er) = openBrowser2(&url){
                eprintln!("Fallbask failed too: \n{}", er);
            }

        }
    }
    let mut version2 = String::new();

    let re = Regex::new(r"https?://(?P<name>[^\.]+)\.org/\?v=(?P<version>[\d\.]+)").unwrap();
   
    if let Some(caps) = re.captures(&version) {
        let name = &caps["name"];                    // "wordpress"
        let version_raw = &caps["version"];          // "6.1"
        let version_formatted = version_raw.replace(".", "+"); // "6+1"

        version2 = format!("{}+{}", name, version_formatted);
        
    } else {
        println!("Not valid url");
    }  

    let url1 = format!("https://www.cve.org/CVERecord/SearchResults?query={}", version2);

    match openBrowser(&url1) {
        Ok(_) => {}
        Err(e) => {
            eprintln!("Trying in another way beacuse of error {}...", e);
            if let Err(er) = openBrowser2(&url1){
                eprintln!("Fallbask failed too: \n{}", er);
            }

        }
    }


    
    Ok(())
}
async fn save_feed(text: &String) -> Result<(),Box<dyn Error>>{
    let mut file = File::create("feed.xml").await?;
    file.write_all(text.as_bytes()).await?;

    Ok(())
}

// if parser fails
fn find(text: &String) -> &str{
    for line in text.lines(){
        if line.contains("generator"){
            return line;
        }
    }
    &"Not found"
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

#[allow(non_snake_case)]
fn openBrowser2(url: &str) -> Result<(),Box<dyn Error>>{
    if webbrowser::open(url).is_ok(){
        Ok(())
    }else {
      Err("Failed to open browser".into()) 
    }
}
#[allow(non_snake_case)]
fn openBrowser(url: &str) -> Result<(), Box<dyn Error>> {
    match OS{
        "windows" => {
            Command::new("cmd")
            .args(["/C", "start", url])
            .spawn()?;
            Ok(())
        }

        "macos" => {
            Command::new("open")
            .arg(url)
            .spawn()?;
            Ok(())
        }

        "linux" => {
            Command::new("xdg-open")
            .arg(url)
            .spawn()?;
            Ok(())
        }

        _ => {eprintln!("Not supported OS");
            Ok(())}
    }
}
