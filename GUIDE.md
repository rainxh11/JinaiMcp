## 🖥️ Usage
Once the Docker container is running, you can use curl to make requests. Here are examples for different response types:

1. 📝 Markdown (bypasses readability processing):
   ```bash
   curl -H "X-Respond-With: markdown" 'http://reader-container:3000/https://google.com'
   ```

2. 🌐 HTML (returns documentElement.outerHTML):
   ```bash
   curl -H "X-Respond-With: html" 'http://reader-container:3000/https://google.com'
   ```

3. 📄 Text (returns document.body.innerText):
   ```bash
   curl -H "X-Respond-With: text" 'http://reader-container:3000/https://google.com'
   ```

4. 📸 Screen-Size Screenshot (returns the URL of the webpage's screenshot):
   ```bash
   curl -H "X-Respond-With: screenshot" 'http://reader-container:3000/https://google.com'
   ```

5.  📸 Full-Page Screenshot (returns the URL of the webpage's screenshot):
   ```bash
   curl -H "X-Respond-With: pageshot" 'http://reader-container:3000/https://google.com'
   ```