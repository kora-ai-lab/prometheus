package com.prometheus.app

import android.os.Bundle
import android.webkit.WebView
import android.webkit.WebViewClient
import android.widget.ProgressBar
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity

class MainActivity : AppCompatActivity() {

    private lateinit var webView: WebView
    private lateinit var statusText: TextView
    private lateinit var progressBar: ProgressBar

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        webView = findViewById(R.id.web_view)
        statusText = findViewById(R.id.status_text)
        progressBar = findViewById(R.id.progress_bar)

        setupWebView()
        loadWebUI()
    }

    private fun setupWebView() {
        webView.settings.javaScriptEnabled = true
        webView.settings.domStorageEnabled = true
        webView.webViewClient = object : WebViewClient() {
            override fun onPageFinished(view: WebView?, url: String?) {
                progressBar.visibility = ProgressBar.GONE
            }
        }
    }

    private fun loadWebUI() {
        progressBar.visibility = ProgressBar.VISIBLE
        val serverUrl = getServerUrl()
        webView.loadUrl(serverUrl)
    }

    private fun getServerUrl(): String {
        val host = intent.getStringExtra("SERVER_HOST") ?: "127.0.0.1"
        val port = intent.getIntExtra("SERVER_PORT", 8080)
        return "http://$host:$port"
    }

    override fun onBackPressed() {
        if (webView.canGoBack()) {
            webView.goBack()
        } else {
            super.onBackPressed()
        }
    }
}
