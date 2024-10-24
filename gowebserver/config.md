```json5
{   
    "http": {
        "server": {
            "listen": 80,

            "location" {
                // Static content
                "/": {
                    "root": "/data/www"
                },
                // "/images/aaa.png" -> "/data/images/aaa.png"
                "/images/": {
                    "root": "/data"
                },

                "~ \.(gif|jpg|png)$": {
                    "root" "/data/images"
                }

                // Proxy Server
                "/": {
                    "proxy_pass": "http://localhost:8080",
                },
                // "/images/aaa.png" -> "http://localhost:8080/images/aaa.png"
                "/images": {
                    "proxy_pass": "http://localhost:8080",
                },

                // "/images/aaa.png" -> "http://localhost:8080/aaa.png"
                "/images/": {
                    "proxy_pass": "http://localhost:8080/",
                }
            }
        }
    }
}
```