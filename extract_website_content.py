import requests
from bs4 import BeautifulSoup
from bs4.element import Comment
from urllib.parse import urljoin, urlparse
import sys
import xml.etree.ElementTree as ET
import time
import json
import uuid

# Replace with your sitemap URL
sitemap_url = 'https://stevenlawton.com/sitemap.xml'

# Web service endpoint to post JSON data
web_service_url = 'http://localhost:8080/api/receive_data'

# Generate a unique extract ID
extract_id = str(uuid.uuid4())

visited = set()
to_visit = []

output_data = []

# Function to parse sitemap and extract URLs
def extract_urls_from_sitemap(sitemap_url):
    try:
        response = requests.get(sitemap_url)
        if response.status_code != 200:
            print(f'Failed to fetch sitemap: {response.status_code}', file=sys.stderr)
            return []
        sitemap_content = response.content
        root = ET.fromstring(sitemap_content)

        urls = []
        # If it's a sitemap index, get all nested sitemap URLs
        for elem in root.iter():
            if elem.tag.endswith("sitemapindex"):
                for sitemap in root.findall(".//{*}sitemap/{*}loc"):
                    sitemap_url = sitemap.text.strip()
                    urls.extend(extract_urls_from_sitemap(sitemap_url))
            elif elem.tag.endswith("loc"):  # Extract all <loc> elements from regular sitemaps
                urls.append(elem.text.strip())

        return urls
    except requests.exceptions.RequestException as e:
        print(f'Failed to fetch sitemap: {e}', file=sys.stderr)
        return []

# Function to check if a tag is visible
def tag_visible(element):
    if element.parent.name in ['style', 'script', 'head', 'title', 'meta', '[document]']:
        return False
    if isinstance(element, Comment):
        return False
    return True

# Extract URLs from the sitemap
to_visit = extract_urls_from_sitemap(sitemap_url)

# Process each URL
while to_visit:
    url = to_visit.pop(0)
    if url in visited:
        continue
    visited.add(url)

    # Measure page load time
    start_time = time.time()
    try:
        response = requests.get(url)
    except requests.exceptions.RequestException as e:
        print(f'Failed to fetch {url}: {e}', file=sys.stderr)
        continue
    load_time = time.time() - start_time

    if response.status_code != 200:
        print(f'Non-200 response for {url}: {response.status_code}', file=sys.stderr)
        continue
    soup = BeautifulSoup(response.content, 'html.parser')

    # Extract page title
    title = soup.title.string if soup.title else 'No Title'

    # Extract meta tags
    meta_tags = {}
    for meta in soup.find_all('meta'):
        name = meta.get('name', '').lower()
        prop = meta.get('property', '').lower()
        if name or prop:
            key = name if name else prop
            meta_tags[key] = meta.get('content', '')

    # Extract canonical tag
    canonical = soup.find('link', rel='canonical')
    canonical_url = canonical['href'] if canonical else 'No Canonical URL'

    # Extract H tags
    h_tags = {}
    for i in range(1, 7):
        h_tag = f'h{i}'
        h_tags[h_tag] = [h.get_text(strip=True) for h in soup.find_all(h_tag)]

    # Extract visible text from the page
    texts = soup.findAll(string=True)
    visible_texts = filter(tag_visible, texts)
    body = u" ".join(t.strip() for t in visible_texts)

    # Calculate word count
    word_count = len(body.split())

    # Extract images and their alt attributes
    images = soup.find_all('img')
    images_info = [{'src': img.get('src'), 'alt': img.get('alt', 'No alt attribute')} for img in images]

    # Extract internal and external links
    internal_links = []
    external_links = []
    for link in soup.find_all('a', href=True):
        href = link['href']
        full_url = urljoin(url, href)
        if urlparse(full_url).netloc == urlparse(sitemap_url).netloc:
            internal_links.append(full_url)
        else:
            external_links.append(full_url)

    # Check for broken links
    broken_links = []
    for link in internal_links + external_links:
        try:
            link_response = requests.head(link, allow_redirects=True)
            if link_response.status_code == 404:
                broken_links.append(link)
        except requests.exceptions.RequestException:
            broken_links.append(link)

    # Extract structured data (JSON-LD)
    structured_data = soup.find_all('script', type='application/ld+json')
    structured_data_content = [sd.get_text(strip=True) for sd in structured_data]

    # Extract robots meta tag
    robots_meta = soup.find('meta', attrs={'name': 'robots'})
    robots_content = robots_meta['content'] if robots_meta else 'No Robots Meta Tag'

    # Title and meta description length
    title_length = len(title) if title else 0
    meta_description = meta_tags.get('description', '')
    meta_description_length = len(meta_description)

    # H1 tag count
    h1_count = len(h_tags.get('h1', []))

    # Create a dictionary for each page's data
    page_data = {
        'extract_id': extract_id,
        'URL': url,
        'Title': title,
        'Title Length': title_length,
        'Meta Description': meta_description,
        'Meta Description Length': meta_description_length,
        'Meta Tags': meta_tags,
        'Canonical URL': canonical_url,
        'H Tags': h_tags,
        'H1 Tag Count': h1_count,
        'Word Count': word_count,
        'Page Load Time (seconds)': round(load_time, 2),
        'Images': images_info,
        'Internal Links': internal_links,
        'External Links': external_links,
        'Broken Links': broken_links,
        'Structured Data': structured_data_content,
        'Robots Meta Tag': robots_content,
        'Content': body
    }

    # Append the page data to the output list
    output_data.append(page_data)

    # Post the page data to the web service
    try:
        post_response = requests.post(web_service_url, json=page_data)
        if post_response.status_code != 201:
            print(f'Failed to post data for {url}: {post_response.status_code} - {post_response.text}', file=sys.stderr)
    except requests.exceptions.RequestException as e:
        print(f'Failed to post data for {url}: {e}', file=sys.stderr)

# Save output to a JSON file
with open('website_content.json', 'w', encoding='utf-8') as f:
    json.dump(output_data, f, indent=4, ensure_ascii=False)

print('Website content has been extracted to website_content.json')
