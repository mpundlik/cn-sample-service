# CN sample service

The CN sample service showcases the extensibility of vpp-agent.

The basic steps to setup a project:

Initialize a tool for management dependencies. This example uses [glide](https://github.com/Masterminds/glide).
We assume that glide is already installed if not, follow the instruction in its README.

```
glide init
```

Modify the content of the `glide.yaml` that defines dependencies of the project.
It is recommanded to pin dependencies to a particular commit id or a tag. The initial content
might look like this:

```yaml
package: cn-sample-service
import:
- package: github.com/ligato/vpp-agent
  version: 9b1e57b07a1dbda7e76f2fb0a7e2f584eb568b92
```

Download initial set of dependencies

```
glide install --strip-vendor
```

Once the initial set of dependencies is downloaded we can move to the writing of [custom plugin](cmd/helloworld).
The common tasks related to project development such as building, updating of dependencies, running of static
analysis and so on... can be automated using Makefile. Take a look at the [Makefile](Makefile) in this repository.