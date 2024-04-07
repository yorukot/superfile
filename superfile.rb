class Superfile < Formula
  desc "Modern and pretty fancy file manager for the terminal"
  homepage "https://github.com/MHNightCat/superfile"
  url "https://github.com/MHNightCat/superfile/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "b91aacb0966dacf92efd27d9bbb4aff7d7b4cdc77168a21880aed1db3e456ffe"

  depends_on "exiftool"

  def install
    bin.install Dir["bin/*"]
  end

  test do
    output = shell_output("#{bin}/spf -v")
    assert_match("superfile version v1.0.0", output)
  end
end
