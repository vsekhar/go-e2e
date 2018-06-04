package ca.viveksekhar.goend_to_end

import android.app.Activity
import android.os.Bundle
import android.content.Intent
import android.net.Proxy.getHost
import android.net.Uri
import android.webkit.*


class MainActivity : Activity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        val myWebView = findViewById(R.id.webView1) as WebView
        val webSettings = myWebView.settings
        webSettings.javaScriptEnabled = true
        myWebView.webViewClient = object: WebViewClient() {
            override fun shouldOverrideUrlLoading(view: WebView, request: WebResourceRequest): Boolean {
                if (request.getUrl().getHost().equals("www.google.com")) {
                    return false
                }
                val intent = Intent(Intent.ACTION_VIEW, request.getUrl())
                startActivity(intent)
                return true
            }
        }
        myWebView.webChromeClient = WebChromeClient()
        myWebView.loadUrl("https://www.google.com")
    }
}

