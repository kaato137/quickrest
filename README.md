# QuickREST

*When you need an API endpoint, fast!*

QuickREST is a convenient tool for quickly mocking API endpoints when you don't have them readily available.

## How It Works

Here's a simple guide on how to use QuickREST:

1. **Define Endpoints**: Begin by creating a YAML file with definitions for your desired endpoints. Specify the address and routes as follows:

```yaml
addr: localhost:8090

routes:

- path: GET /api/1/articles/{id}
  status: 200
  body: |
    {
        "items": [
            {
                "id": {id},
                "headline": "QuickREST is a hot new thing"
            }
        ]
    }
```

2. **Run the CLI**: Execute the QuickREST CLI against the configuration file:

```bash
quickrest -c quickrest.yml
```

That's it! You now have a basic REST API server up and running, ready to serve your mocked endpoints.

## Installation

Just clone this repository and run:

```bash
go install ./...
```

This would compile and place `quickrest` binary in your `$GOPATH/bin`.

## Key Features

- **Rapid Mocking**: Quickly create mocked API endpoints.
- **YAML Configuration**: Define endpoints easily using YAML syntax.
- **Simple CLI**: Run QuickREST with a straightforward command-line interface.

## Templating

QuickREST supports two kinds of templating: basic and JavaScript. Basic templating allows you to insert URL parameters inside the response body like so:

```yaml
routes:

- path: GET /hello/{name}
  body: |
    Hello, {name}!
```

The JavaScript templating feature enables you to create complex responses:

```yaml
routes:

- path: GET /api/1/users/{user_id}/notifications
  body_js: |
    var result = [];

    for (var i = 0; i < 5; i++) {
        result.push({
            "id": parseInt(user_id) + i,
            "user_id": parseInt(user_id),
            "description": "You got mail!"
        });
    }

    result;
```
In the JavaScript templating example, you can utilize JavaScript code to dynamically generate the response body. This allows for more flexibility in crafting responses based on dynamic data or complex logic.

## Get Started

To get started with QuickREST, simply clone this repository and follow the instructions above to create and run your mocked endpoints.

## Contributing

If you find issues or have ideas to improve QuickREST, feel free to contribute by forking this repository, making your changes, and submitting a pull request. We welcome any contributions that can enhance the functionality and usability of QuickREST.

## License

QuickREST is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
Designed and maintained by [Andrei Kuzmin](https://github.com/kaato137) - [Contact Me](mailto:kaato361@gmail.com)