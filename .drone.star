# workaround to render locally since you cant pass repo.branch to the cli
def repo_branch(ctx):
    return getattr(ctx.repo, "branch", "main")


def new_pipeline(name, arch, **kwargs):
    pipeline = {
        "kind": "pipeline",
        "name": name,
        "platform": {
            "arch": arch,
        },
        "steps": [],
    }

    pipeline.update(kwargs)

    return pipeline


def pipeline_test(ctx):
    cache_volume = {"name": "cache", "temp": {}}
    cache_mount = {"name": "cache", "path": "/go"}

    # licensed-go image only supports amd64
    p = new_pipeline(
        name="test",
        arch="amd64",
        trigger={"branch": repo_branch(ctx)},
        volumes=[cache_volume],
        workspace={"path": "/go/src/github.com/{}".format(ctx.repo.slug)},
        steps=[
            {
                "commands": ["make test"],
                "image": "golangci/golangci-lint",
                "name": "test",
                "volumes": [cache_mount],
            },
            {
                "commands": ["licensed cache", "licensed status"],
                "image": "public.ecr.aws/kanopy/licensed-go",
                "name": "license-check",
            },
        ],
    )

    return p


def main(ctx):
    pipelines = [pipeline_test(ctx)]

    return pipelines
