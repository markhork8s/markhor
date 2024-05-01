package clientset

import (
	"testing"

	"k8s.io/apimachinery/pkg/version"
)

func TestIsCompatible(t *testing.T) {
	t.Run("A known compatible version is compatible", func(t *testing.T) {
		v := version.Info{
			Major: "1",
			Minor: "28",
		}
		if !IsCompatible(&v) {
			t.Fatal("Should have been compatible")
		}
	})
	t.Run("An invalid version is incompatible", func(t *testing.T) {
		v := version.Info{
			Major: "-1",
			Minor: "28",
		}
		if IsCompatible(&v) {
			t.Fatal("Should have been incompatible")
		}
	})
	t.Run("An unknown version is incompatible", func(t *testing.T) {
		v := version.Info{
			Major: "9999999999999999",
			Minor: "28",
		}
		if IsCompatible(&v) {
			t.Fatal("Should have been incompatible")
		}
	})
}
