package myproject;

public class Environment {
    public enum Size {SMALL, MEDIUM, LARGE}

    public enum Kind {DEV, STAGING, PROD}

    private String name;
    private Kind kind;
    private Size size;
    private String gitOpsURL;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Kind getKind() {
        return kind;
    }

    public void setKind(Kind kind) {
        this.kind = kind;
    }

    public Size getSize() {
        return size;
    }

    public void setSize(Size size) {
        this.size = size;
    }

    public String getGitOpsURL() {
        return gitOpsURL;
    }

    public void setGitOpsURL(String gitOpsURL) {
        this.gitOpsURL = gitOpsURL;
    }
}
