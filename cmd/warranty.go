package cmd

import "github.com/spf13/cobra"

// warrantyText mirrors the disclaimer paragraph in LICENSE verbatim, so
// there's a single source of truth for the actual legal wording.
const warrantyText = `NO WARRANTY

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

See LICENSE for the full MIT License text.`

var warrantyCmd = &cobra.Command{
	Use:   "warranty",
	Short: "Show the no-warranty disclaimer",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println(warrantyText)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(warrantyCmd)
}
