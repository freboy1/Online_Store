package auth

import (
    "testing"
    "github.com/tebeka/selenium"
    "time"
    "fmt"
)

func TestLoginPage(t *testing.T) {
    // Setup WebDriver capabilities for Microsoft Edge
    caps := selenium.Capabilities{"browserName": "MicrosoftEdge"}
    edgeCaps := selenium.Capabilities{
        "ms:edgeOptions": map[string]interface{}{
            "args": []string{"--disable-gpu"},
        },
    }

    // Merge the Edge-specific capabilities with the base capabilities
    for key, value := range edgeCaps {
        caps[key] = value
    }

    // Connect to Edge WebDriver running on port 5555
    wd, err := selenium.NewRemote(caps, "http://127.0.0.1:52090")
    if err != nil {
        t.Fatal(err)
    }
    defer wd.Quit()

    // Navigate to the login page
    err = wd.Get("http://127.0.0.1:8080/login")
    if err != nil {
        t.Fatal(err)
    }

    // Find and fill the email input
    emailInput, err := wd.FindElement(selenium.ByID, "email")
    if err != nil {
        t.Fatal(err)
    }
    err = emailInput.SendKeys("dautovalisher33@gmail.com")
    if err != nil {
        t.Fatal(err)
    }

    // Find and fill the password input
    passwordInput, err := wd.FindElement(selenium.ByID, "password")
    if err != nil {
        t.Fatal(err)
    }
    err = passwordInput.SendKeys("123")
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(5 * time.Second)
    // Find the submit button using the correct ID (login)
    submitButton, err := wd.FindElement(selenium.ByID, "login")
    if err != nil {
        t.Fatal(err)
    }

    // Use JavaScript to click the submit button
    jsScript := "arguments[0].click();"
    _, err = wd.ExecuteScript(jsScript, []interface{}{submitButton}) // Capture both result and error
    if err != nil {
        t.Fatal("Failed to execute JavaScript click:", err)
    }

    // Allow time for form submission and page load
    time.Sleep(5 * time.Second)

    // Get the cookies after login
    cookies, err := wd.GetCookies()
    if err != nil {
        t.Fatal("Failed to retrieve cookies:", err)
    }

    // Print cookies for debugging (optional)
    if len(cookies) == 0 {
        t.Error("No cookies found after login.")
    } else {
        fmt.Println("Cookies after login:")
        for _, cookie := range cookies {
            fmt.Printf("Cookie: Name: %s, Value: %s\n", cookie.Name, cookie.Value)
        }
    }

    // Verify the resulting URL
    currentURL, err := wd.CurrentURL()
    if err != nil {
        t.Fatal(err)
    }

    // Adjust the expected URL based on your application’s behavior
    if currentURL != "http://127.0.0.1:8080/products" {
        t.Errorf("Expected URL to be 'http://127.0.0.1:8080/products', but got '%s'", currentURL)
    }
}
