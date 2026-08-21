package domain

import "testing"

func TestSide(t *testing.T) {
	tests := []struct {
		name       string
		side       Side
		wantIsBuy  bool
		wantIsSell bool
		wantString string
	}{
		{
			name:       "buy side",
			side:       SideBuy,
			wantIsBuy:  true,
			wantIsSell: false,
			wantString: "BUY",
		},
		{
			name:       "sell side",
			side:       SideSell,
			wantIsBuy:  false,
			wantIsSell: true,
			wantString: "SELL",
		},
		{
			name:       "unknown side",
			side:       SideUnknown,
			wantIsBuy:  false,
			wantIsSell: false,
			wantString: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.side.IsBuy(); got != tt.wantIsBuy {
				t.Errorf("IsBuy() = %v, want %v", got, tt.wantIsBuy)
			}
			if got := tt.side.IsSell(); got != tt.wantIsSell {
				t.Errorf("IsSell() = %v, want %v", got, tt.wantIsSell)
			}
			if got := tt.side.String(); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
		})
	}
}
