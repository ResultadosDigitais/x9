
<p  align="center">
<img  width="150"  src="docs/img/gopherx9.jpg"  alt="X9"  title="X9"  />
</p>

## X9: Static analysis in real-time

X9 is a simple static analysis tool that helps find potential sensitive data leaks. It integrates with Github webhooks to receive Pull Request events and can send the results to slack.

### How it works

Developers are constantly creating new features making it hard for security teams to manually analyze every piece of code. The purpose of this project is to bring security automation to find potential sensitive data exposure in the development process.

Let's suppose that we have a development team working on a new feature. At some point, they make a new Pull Request to allow code reviews. Github will receive it and it will notify X9 by sending the Pull Request data to X9 `/events` route.

X9 has many workers, which can be configured to deal with a great number of requests. It will send this event to one of the workers that are responsible for performing the tests. First, the worker will clone the repository in a temporary folder. Then it will analyze every single file searching for a match on one of the configured signatures. When it finds a match, it logs to `stdout` and sends it to slack through a configured webhook.

<p  align="center">
<img  width="800"  src="docs/img/diagram.png"  alt="diagram"  title="diagram"  />
</p>

### Customize


#### Signatures

Signatures are regular expressions that X9 will use to found potential sensitive data in your code. You can disable or add new signatures by editing `config.yaml` file.

Signatures have some fields that you need to set:
  

-  **part**
part is the resource type that x9 will apply the regular expression. It has to be set with one of the following possible values:
   -  *extension*: matches files extensions
   -  *path*: matches full path values
   -  *filename*: matches filename only
   -  *contents*: matches values inside the files

-  **regex**
The regular expression that will be used to match values in the respective resource.

-  **name**
Vulnerability name that will be displayed on the report.