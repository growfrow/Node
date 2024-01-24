package paragraph_test

import (
	"context"
	"testing"

	"github.com/rss3-network/serving-node/config"
	source "github.com/rss3-network/serving-node/internal/engine/source/arweave"
	worker "github.com/rss3-network/serving-node/internal/engine/worker/contract/paragraph"
	"github.com/rss3-network/serving-node/provider/arweave"
	"github.com/rss3-network/serving-node/schema"
	"github.com/rss3-network/serving-node/schema/filter"
	"github.com/rss3-network/serving-node/schema/metadata"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestWorker_Arweave(t *testing.T) {
	t.Parallel()

	type arguments struct {
		task   *source.Task
		config *config.Module
	}

	testcases := []struct {
		name      string
		arguments arguments
		want      *schema.Feed
		wantError require.ErrorAssertionFunc
	}{
		{
			name: "Paragraph Post",
			arguments: arguments{
				task: &source.Task{
					Network: filter.NetworkArweave,
					Block: arweave.Block{
						Timestamp: 1697091466,
					},
					Transaction: arweave.Transaction{
						ID:       "Sz5fY8Loj67fWxLQv98r5U5-h2aIA5x4FMsAVP1N2ig",
						Reward:   "212017846",
						Quantity: "0",
						Owner:    "rsfrA1_2H7Pb4kRtHj6EryEELG1sksd-1xGbAWJJqgCIJs9dQYL2C7afuCFX-pryKFpU3ZLssERyObt-BiDwWA3vSHAFljt0CbCBZRKqWKWeEXXdoBLR_Vf8724P14YqRubW7a0n6UaZKsJsxah35yPCANnw9QbnHJouTlNyky41ZnbBClRlYWr1_PkEMvFEsQcqIE5J8jcgJlaTNtiOi7ruvRP3z-NtqufuJFFq3_4hrL6ICpbJnZBgZuX33tr6YvCrYExtFmd8wJoL4s6MSioKYSWYk60ngr8EgUHotS1lzPemWRhY9wjbrg3wh00sCO44wv5CmE2Ke-EoKZYKrUu8g5z2MlPwOnNxBj67wzrSRzkpaVvbEWqneEwG_UcDxKU_SLeJ0_qGLNkQjgqjhfAAEdivsfV0Fz3hNRmVu2ae84QtoPQyvvcr2JLe-bTjbGvna_C52fR7-p9sp-MlZnL8vPnKfPZrTvfCOd935O2_CdiyzvOA35jQKQhe5UhqwH0hoYdplE2DHRN6MR42n-8nq3vqxp7Y34l-aUxnRIHBquMFbfH4KKn8N322_e_6nAwImjp_DziPhz5xOyQJgZOzCBTFuQrbaHkGbQ6ou814fyAUDJlA3S5-WKtsD8Jk1AMg0YmIdFUgCVUwwepoAgK1UPAxpq64GouKmnqjI58",
						Data:     "eyJ1Z2MiOnRydWUsImFyY2hpdmVkIjpmYWxzZSwiaWQiOiJZZUNpZENXME1qOVFyajZoSjhnWCIsImFjY2Vzc1Jlc3RyaWN0aW9uIjpudWxsLCJhdXRob3JzIjpbIm1Ld0R6V0VIUURvZjhIcHc5cHlyIl0sInRpdGxlIjoiWXUgVGVzdCBQb3N0Iiwic3VidGl0bGUiOiJQb3N0IiwiY3JlYXRlZEF0IjoxNjk3MDkxMTQyNzI3LCJjYXRlZ29yaWVzIjpbImRhdGEiXSwianNvbiI6IntcInR5cGVcIjpcImRvY1wiLFwiY29udGVudFwiOlt7XCJ0eXBlXCI6XCJoZWFkaW5nXCIsXCJhdHRyc1wiOntcInRleHRBbGlnblwiOlwibGVmdFwiLFwibGV2ZWxcIjoxfSxcImNvbnRlbnRcIjpbe1widHlwZVwiOlwidGV4dFwiLFwidGV4dFwiOlwidGhhdCdzIG15IGZpcnN0IHBvc3QgaGVyZVwifV19LHtcInR5cGVcIjpcImZpZ3VyZVwiLFwiYXR0cnNcIjp7XCJzcmNcIjpudWxsLFwiZmlsZVwiOm51bGwsXCJhbHRcIjpudWxsLFwidGl0bGVcIjpudWxsLFwiYmx1ckRhdGFVUkxcIjpudWxsLFwiZmxvYXRcIjpcIm5vbmVcIixcIndpZHRoXCI6bnVsbH0sXCJjb250ZW50XCI6W3tcInR5cGVcIjpcImltYWdlXCIsXCJhdHRyc1wiOntcInNyY1wiOlwiaHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL3BhcHlydXNfaW1hZ2VzLzk3ZjM3YmU3NDIyNTJhMmRhNTBhYjlhYzBmM2E0ODUxLmpwZ1wiLFwiYWx0XCI6bnVsbCxcInRpdGxlXCI6bnVsbCxcImJsdXJkYXRhdXJsXCI6XCJkYXRhOmltYWdlL3BuZztiYXNlNjQsaVZCT1J3MEtHZ29BQUFBTlNVaEVVZ0FBQUNBQUFBQVZDQUlBQUFDb3IzdTlBQUFBQ1hCSVdYTUFBQXNUQUFBTEV3RUFtcHdZQUFBSDYwbEVRVlI0bkNYQmVWVFNlUUlBOEcrenpiemROelB2Ylc5MzUvaGpwM2JxZFV4VFRiNXlLaXQxSisyd083TkRMY3NvNzlBU3NCUVZ4YlBTRkFVRlJmMVpHb28vL1lrZ2Nrb2d2MFFVa0t1VVE4VVFqN1FtdFZMNzdSLzcrWUN2UTJrQlZ6S1RMbDRtWDc1MHBoeFpyVmo0RXYxenRXSm10V0x5TC9LeE5STGJWNktYNjRYVy9TS2psOGpvS3pINUM3WDdCWnJEUXMwQmtkNUhhTnpOMTIzaDZmNk9ESzVDektESkFCb05nR01GMVdidytBWElVWUJVSWZqZDIrZTcwTHp0Y1JWSjEyT0NjRWtncS8wTGdjMWI1ZnFyeXJNSmRlOUd4LzZyc2gxQmJmc1ZWb0ptWUhCWWlmUnJjQXA5Q0dvK3BUS2ZVNWo4cEJadm9UNUFxQXVXbU5aM21FR3JlVlhyMEdaWUM1Z0tVQ29CdWUzZzZKYS9IZGoyMDg4blkwN2ppMjVHeFFRbXBYSDRzRTZ2VkEyYnV4d094T0dBbkk0VTNURFhZcHgvamI1M3FCeldBVVJsaUd3YmFEUnFQSytWb2xlYUZxT21ycSszQUIzSTZyVVFrSzdIVXVVcmc5clFxMEpsVWdUaGc2TGduOHRDMW1VZldaTisyWmRka01wa1VIbGlUcU5HeVJzeXduWWIwekFTeVhQNnMrMEVydjAyN0R6VE9QSXJjeGlVRElGVVBjanBvM1dwcDUxeWJFYXg2Tzc1YzFTbGt0YUdoaDRqa2VMMXluWmVLNjBmWVNycm1NQ0kvOUYyYjVNK2NTMFU2emNHcHp0RU9haVlWU2ZrazF0NmNFLzBvWkFsSHJibUt1eVArMGJwcGxHVzFWVmpjVUJXTjhQa3ZpNTg3UTg1b3A1YW0yU0RkcXRjcklTU1UyTHcrT2kwOU9TN0JjSEJVWHNmRVJMMmJ0OEJuQmxlN2x3ZkMzRkxLekZJQzVHNWowZ3dvNUJYdzVUeUVFT2Z6T1BxVzV3elkrOHRLL1BXcFFYN3U3ZjJwZmVXOFlrZW5SMTk2UnBBaCtSYXE2aEpqWEQ3bXlJcHVQMGhJVDR4OTMyaXduY0hyNjNOUzl5NmVRLzQ0bHZnb2g3d0ZQcU9VWGJ4eWFjZFhNcGtDL2xEOTBQTVdJblpPZGlFQ0p0Vkw4MTB6STUzWVcvMWI2Y0hYcnBlREk3MHFvYmxBcE8wd3lTVERDblVEcVhSb3hRTXRJZkZYdmM3RVhLTHVNODNZbE40aEc4SkdmL2RmM1p0MkxnRG1Nbitsa3cvUFhGdksra013c2hLSWFSa1UzTXpIdEl5YUZEZUU0UUJ0MVUyRmxLZ3lseUIrb0hjOUFpMVB0SU1seHZjSEpPclhqOVViUmpLYWU5TktZS3VFa3FDcnQyL2MvUFM3UXQvbkRvV21KWWN5eVplTzcxMVk4Q1dqWUNYUzVROEpQTHpFL2tGZUNubklZTk5ld0t6T2Z3RzRmUDJWbzFjcE85UkRCazFveTlmak50ZVRMcjZwOTJHMmNtK1VVVHBRQWRuUGFoTkYwdEt2Wm9RR1o0YWRqZDNSem9CcEJOQVJNaS84bkxPc0hCQmV3RFlCZ0JZMHNBck9uakYwTHlrWm5HYktzUHFPQ1ErLzRGY1NrZTdXV1lkMTJia09peXd5eTZZY01yZXVHV3pVL0s1S2RocGhGODdCTk51c2ZzVk90RnZuZW1jbXl4ODQ0aUdvY0NzQjhEL05NQW43cTY4ZWNJTGdPMWZBYkNvYXZpb2JWNDJjSmZRbW5KYTBZbHNXdXl6K2hqTzAzdTgxbnM4aENRWEZodFVOTE9tWUxBL1k5QlU2VEEvR3pjMnVXMDhqMVhtTWZTNmRDLzZlVk5PR2VidUhETW5JNXpmb3BQWDdEb0dyb1I3VXk0YzNRREFMd0NBdCtLYUJWWERSeDEzRVcwb3BxVnRKOFlmWjJiNjAxTndVRlpFWTJsd1MwVVN5bVZaQmR4UnVYaGMzVDlyTU0vMGowLzF6RThJVnp4dDJBSnZoRTloRVU2bVoyNklUd1ZCbDhEQm8rQ1huZUI2NkQ3U3VTUC9CbUF0QUdBS3FaaVQxaXlvRzJka05aTDZiSG9sa1FzUk9yc0twU3FHVXB4bU05ZWJ0ZldlSVhobW1JTjkxTTUyc3o2cFdaOHN6OTY3K0N1VHZIZEcydndJVDhrbEtubnBwWVFMYmNWSnlrYnFzOHdyMGxKODBwbkQzd0x3VHdDQUEzbzBBZFBmeWRpamdnb3pYT3hxejU5dlQzNkRNdWZOenhiVm1TdU9acU9NN3RKRHoyVVpIU25oOGwrL212Y0NrMGZYTGJia2Y1aVdPS1U1Yy9hT2FVM3QyMWM4SlpRMjFrMmZNM0NrTkx5V1RZbzdlbkFWQUY4REFBYUtLYmFhQjU0V21vMWI0aEl4M3FEc09RTjNjYmpqNDRob2VVS3hNdG1EemZaaTJMQ3BQb2Q0UG9FYmUxd2UrSTJHbmU5a1Z5MnJxckRsVjlpTWVubWkrNE5UK0hsY3ZHQ0I1L29iWFlwYVYwZEozUEU5NFA5RTZRUzBnR1RJang1aVowK0pheGEwbkFVOWQ4R01ZSy9WMkJzZE5xUEY1dlNmNS9TU3BMQ2laR3JKb2U4OTE5WXZTbGwyVlBmcEtSVmJ0bUd6L2RoTUwrWldMVHNrbjE1MmZyTHduT0pxVXkybE92YlU3MnQvOU4rMEFVQnhONUQ3Y2ZVRnFWUWFoVWkrUWFiY3lNbkZTVHZMUkgxTlRDVHpDUys5V1pKRks4RTFmd25tTm9ENUg4RG5uZUR6K2U4eHZEOUdQb0x3aVBWaUNzUWpWOFAzR0kzRW9xcmJWRm8wSVQwNE1pNGdJZWxLTmVGR0J5a2FGQkNTSSs3YzhVc203a21JOFQ1NzRHUmtVSFh0UTZZWWltWVJJOHFpb3lyamNFeDhWT3h4Nk1TMjhiQk5TOUZlUy9nRGk4bitTK1c0RC9UNHlySVlRa1BLdlRvaXNZNllWSDNuMXVPbzhOeUljOFRnUTdjT0g3cDFMaVFtSXZiNlJlQ3ovN2VvK0lqaUl1cGRZZ3lUWFFJMU1PanNQQ290dWVwcFJnT1NEd3NMK1hKYVYwK1pTRnZWcVdVS05PV0NQcnFndDd4em9FcW9ZWXBVcFh4cEljelBib0F6NUJDSlRZdlBvdDRzU0xsS1RieDhHWGYyUkxEZlB0K3Q0QjgvZkVQT1RwT3BubWNVWlhQNU1DeEdhbUIyU25GU1VWMUdSWE11MUZIQWxaYTBLV2dDZFlXNHQxS3NxWkJvcTZSOUxMR0cxYWxtZER5bndkSWlUa2RlTlNlVHlTUlFDK09KR2JjU0NWZHVScDgvR3hyNHgyRXZINS9Od0h2THVrQ2Y3ZUhuQWhKeEY0bHhZWlRFcTR5c0JPakJuZWF5K3p4MnBwUlRnUEpwR2lGZEw2MHlkZGNNZGtPRHN1ckI3aHFUakszbGxRNEk2Ym91VmgrZnFlTFMydWhwMVpseHBZVElvcnZYc3VQQzhCY0RMaHpjZkdUblQvOERGT0s5a1ZZSU8rVUFBQUFBU1VWT1JLNUNZSUk9XCIsXCJuZXh0aGVpZ2h0XCI6MTM1NyxcIm5leHR3aWR0aFwiOjIwNDh9fSx7XCJ0eXBlXCI6XCJmaWdjYXB0aW9uXCJ9XX0se1widHlwZVwiOlwicGFyYWdyYXBoXCIsXCJhdHRyc1wiOntcInRleHRBbGlnblwiOlwibGVmdFwifSxcImNvbnRlbnRcIjpbe1widHlwZVwiOlwidGV4dFwiLFwidGV4dFwiOlwidGhhdCdzIHRoZSBjb250ZW50XCJ9XX1dfSIsInN0YXRpY0h0bWwiOiI8aDE-dGhhdCdzIG15IGZpcnN0IHBvc3QgaGVyZTwvaDE-PGZpZ3VyZSBmbG9hdD1cIm5vbmVcIiBkYXRhLXR5cGU9XCJmaWd1cmVcIiBjbGFzcz1cImltZy1jZW50ZXJcIiBzdHlsZT1cIm1heC13aWR0aDogbnVsbDtcIj48aW1nIHNyYz1cImh0dHBzOi8vc3RvcmFnZS5nb29nbGVhcGlzLmNvbS9wYXB5cnVzX2ltYWdlcy85N2YzN2JlNzQyMjUyYTJkYTUwYWI5YWMwZjNhNDg1MS5qcGdcIiBibHVyZGF0YXVybD1cImRhdGE6aW1hZ2UvcG5nO2Jhc2U2NCxpVkJPUncwS0dnb0FBQUFOU1VoRVVnQUFBQ0FBQUFBVkNBSUFBQUNvcjN1OUFBQUFDWEJJV1hNQUFBc1RBQUFMRXdFQW1wd1lBQUFINjBsRVFWUjRuQ1hCZVZUU2VRSUE4Ryt6emJ6ZE56UHZiVzkzNS9oanAzYnFkVXhUVGI1eUtpdDFKKzJ3TzdORExjc283OUFTc0JRVnhiUFNGQVVGUmYxWkdvby8vWWtnY2tvZ3YwUVVrS3VVUThVUWo3UW10Vkw3N1IvNytZQ3ZRMmtCVnpLVExsNG1YNzUwcGh4WnJWajRFdjF6dFdKbXRXTHlML0t4TlJMYlY2S1g2NFhXL1NLamw4am9Lekg1QzdYN0JackRRczBCa2Q1SGFOek4xMjNoNmY2T0RLNUN6S0RKQUJvTmdHTUYxV2J3K0FYSVVZQlVJZmpkMitlNzBMenRjUlZKMTJPQ2NFa2dxLzBMZ2MxYjVmcXJ5ck1KZGU5R3gvNnJzaDFCYmZzVlZvSm1ZSEJZaWZScmNBcDlDR28rcFRLZlU1ajhwQlp2b1Q1QXFBdVdtTlozbUVHcmVWWHIwR1pZQzVnS1VDb0J1ZTNnNkphL0hkajIwODhuWTA3amkyNUd4UVFtcFhINHNFNnZWQTJidXh3T3hPR0FuSTRVM1REWFlweC9qYjUzcUJ6V0FVUmxpR3diYURScVBLK1ZvbGVhRnFPbXJxKzNBQjNJNnJVUWtLN0hVdVVyZzlyUXEwSmxVZ1RoZzZMZ244dEMxbVVmV1pOKzJaZGRrTXBrVUhsaVRxTkd5UnN5d25ZYjB6QVN5WFA2cyswRXJ2MDI3RHpUT1BJcmN4aVVESUZVUGNqcG8zV3BwNTF5YkVheDZPNzVjMVNsa3RhR2hoNGprZUwxeW5aZUs2MGZZU3JybU1DSS85RjJiNU0rY1MwVTZ6Y0dwenRFT2FpWVZTZmtrMXQ2Y0UvMG9aQWxIcmJtS3V5UCswYnBwbEdXMVZWamNVQldOOFBrdmk1ODdRODVvcDVhbTJTRGRxdGNySVNTVTJMdytPaTA5T1M3QmNIQlVYc2ZFUkwyYnQ4Qm5CbGU3bHdmQzNGTEt6RklDNUc1ajBnd281Qlh3NVR5RUVPZnpPUHFXNXd6WSs4dEsvUFdwUVg3dTdmMnBmZVc4WWtlblIxOTZScEFoK1JhcTZoSmpYRDdteUlwdVAwaElUNHg5MzJpd25jSHI2M05TOXk2ZVEvNDRsdmdvaDd3RlBxT1VYYnh5YWNkWE1wa0MvbEQ5MFBNV0luWk9kaUVDSnRWTDgxMHpJNTNZVy8xYjZjSFhycGVESTcwcW9ibEFwTzB3eVNURENuVURxWFJveFFNdElmRlh2YzdFWEtMdU04M1lsTjRoRzhKR2YvZGYzWnQyTGdEbU1uK2xrdy9QWEZ2SytrTXdzaEtJYVJrVTNNekh0SXlhRkRlRTRRQnQxVTJGbEtneWx5QitvSGM5QWkxUHRJTWx4dmNISk9yWGo5VWJSakthZTlOS1lLdUVrcUNydDIvYy9QUzdRdC9uRG9XbUpZY3l5WmVPNzExWThDV2pZQ1hTNVE4SlBMekUva0ZlQ25uSVlOTmV3S3pPZndHNGZQMlZvMWNwTzlSREJrMW95OWZqTnRlVExyNnA5MkcyY20rVVVUcFFBZG5QYWhORjB0S3Zab1FHWjRhZGpkM1J6b0JwQk5BUk1pLzhuTE9zSEJCZXdEWUJnQlkwc0FyT25qRjBMeWtabkdiS3NQcU9DUSsvNEZjU2tlN1dXWWQxMmJrT2l5d3l5NlljTXJldUdXelUvSzVLZGhwaEY4N0JOTnVzZnNWT3RGdm5lbWNteXg4NDRpR29jQ3NCOEQvTk1BbjdxNjhlY0lMZ08xZkFiQ29hdmlvYlY0MmNKZlFtbkphMFlsc1d1eXoraGpPMDN1ODFuczhoQ1FYRmh0VU5MT21ZTEEvWTlCVTZUQS9HemMydVcwOGoxWG1NZlM2ZEMvNmVWTk9HZWJ1SERNbkk1emZvcFBYN0RvR3JvUjdVeTRjM1FEQUx3Q0F0K0thQlZYRFJ4MTNFVzBvcHFWdEo4WWZaMmI2MDFOd1VGWkVZMmx3UzBVU3ltVlpCZHhSdVhoYzNUOXJNTS8wajAvMXpFOElWenh0MkFKdmhFOWhFVTZtWjI2SVR3VkJsOERCbytDWG5lQjY2RDdTdVNQL0JtQXRBR0FLcVppVDFpeW9HMmRrTlpMNmJIb2xrUXNST3JzS3BTcUdVcHhtTTllYnRmV2VJWGhtbUlOOTFNNTJzejZwV1o4c3o5NjcrQ3VUdkhkRzJ2d0lUOGtsS25ucHBZUUxiY1ZKeWticXM4d3IwbEo4MHBuRDN3THdUd0NBQTNvMEFkUGZ5ZGlqZ2dvelhPeHF6NTl2VDM2RE11Zk56eGJWbVN1T1pxT003dEpEejJVWkhTbmg4bCsvbXZjQ2swZlhMYmJrZjVpV09LVTVjL2FPYVUzdDIxYzhKWlEyMWsyZk0zQ2tOTHlXVFlvN2VuQVZBRjhEQUFhS0tiYWFCNTRXbW8xYjRoSXgzcURzT1FOM2Niamo0NGhvZVVLeE10bUR6ZlppMkxDcFBvZDRQb0ViZTF3ZStJMkduZTlrVnkycnFyRGxWOWlNZW5taSs0TlQrSGxjdkdDQjUvb2JYWXBhVjBkSjNQRTk0UDlFNlFTMGdHVElqeDVpWjArSmF4YTBuQVU5ZDhHTVlLL1YyQnNkTnFQRjV2U2Y1L1NTcExDaVpHckpvZTg5MTlZdlNsbDJWUGZwS1JWYnRtR3ovZGhNTCtaV0xUc2tuMTUyZnJMd25PSnFVeTJsT3ZiVTcydC85TiswQVVCeE41RDdjZlVGcVZRYWhVaStRYWJjeU1uRlNUdkxSSDFOVENUekNTKzlXWkpGSzhFMWZ3bm1Ob0Q1SDhEbm5lRHorZTh4dkQ5R1BvTHdpUFZpQ3NRalY4UDNHSTNFb3FyYlZGbzBJVDA0TWk0Z0llbEtOZUZHQnlrYUZCQ1NJKzdjOFVzbTdrbUk4VDU3NEdSa1VIWHRRNllZaW1ZUkk4cWlveXJqY0V4OFZPeHg2TVMyOGJCTlM5RmVTL2dEaThuK1MrVzREL1Q0eXJJWVFrUEt2VG9pc1k2WVZIM24xdU9vOE55SWM4VGdRN2NPSDdwMUxpUW1JdmI2UmVDei83ZW8rSWppSXVwZFlneVRYUUkxTU9qc1BDb3R1ZXBwUmdPU0R3c0wrWEphVjArWlNGdlZxV1VLTk9XQ1BycWd0N3h6b0Vxb1lZcFVwWHhwSWN6UGJvQXo1QkNKVFl2UG90NHNTTGxLVGJ4OEdYZjJSTERmUHQrdDRCOC9mRVBPVHBPcG5tY1VaWFA1TUN4R2FtQjJTbkZTVVYxR1JYTXUxRkhBbFphMEtXZ0NkWVc0dDFLc3FaQm9xNlI5TExHRzFhbG1kRHlud2RJaVRrZGVOU2VUeVNSUUMrT0pHYmNTQ1ZkdVJwOC9HeHI0eDJFdkg1L053SHZMdWtDZjdlSG5BaEp4RjRseFlaVEVxNHlzQk9qQm5lYXkrengycHBSVGdQSnBHaUZkTDYweWRkY01ka09Ec3VyQjdocVRqSzNsbFE0STZib3VWaCtmcWVMUzJ1aHAxWmx4cFlUSW9ydlhzdVBDOEJjRExoemNmR1RuVC84REZPSzlrVllJTytVQUFBQUFTVVZPUks1Q1lJST1cIiBuZXh0aGVpZ2h0PVwiMTM1N1wiIG5leHR3aWR0aD1cIjIwNDhcIiBjbGFzcz1cImltYWdlLW5vZGUgZW1iZWRcIj48ZmlnY2FwdGlvbiBodG1sYXR0cmlidXRlcz1cIltvYmplY3QgT2JqZWN0XVwiIGNsYXNzPVwiaGlkZS1maWdjYXB0aW9uXCI-PC9maWdjYXB0aW9uPjwvZmlndXJlPjxwPnRoYXQncyB0aGUgY29udGVudDwvcD4iLCJzZW5kWE1UUCI6ZmFsc2UsImRvbnRQdWJsaXNoT25saW5lIjpmYWxzZSwic2VuZE5ld3NsZXR0ZXIiOmZhbHNlLCJzdG9yZU9uQXJ3ZWF2ZSI6dHJ1ZSwicG9zdF9wcmV2aWV3IjoidGhhdCdzIG15IGZpcnN0IHBvc3QgaGVyZXRoYXQncyB0aGUgY29udGVudCIsImNvdmVyX2ltZyI6eyJpbWciOnsic3JjIjoiaHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL3BhcHlydXNfaW1hZ2VzL2Y5ZDk1ZTZlZGVkNGQwMGE0ZDc1MmY0ZDAwNGMxYzI5LmpwZyIsIndpZHRoIjoyMDQ4LCJoZWlnaHQiOjEzNTd9LCJpc0hlcm8iOnRydWUsImJhc2U2NCI6ImRhdGE6aW1hZ2UvcG5nO2Jhc2U2NCxpVkJPUncwS0dnb0FBQUFOU1VoRVVnQUFBQ0FBQUFBVkNBSUFBQUNvcjN1OUFBQUFDWEJJV1hNQUFBc1RBQUFMRXdFQW1wd1lBQUFINjBsRVFWUjRuQ1hCZVZUU2VRSUE4Ryt6emJ6ZE56UHZiVzkzNS9oanAzYnFkVXhUVGI1eUtpdDFKKzJ3TzdORExjc283OUFTc0JRVnhiUFNGQVVGUmYxWkdvby8vWWtnY2tvZ3YwUVVrS3VVUThVUWo3UW10Vkw3N1IvNytZQ3ZRMmtCVnpLVExsNG1YNzUwcGh4WnJWajRFdjF6dFdKbXRXTHlML0t4TlJMYlY2S1g2NFhXL1NLamw4am9Lekg1QzdYN0JackRRczBCa2Q1SGFOek4xMjNoNmY2T0RLNUN6S0RKQUJvTmdHTUYxV2J3K0FYSVVZQlVJZmpkMitlNzBMenRjUlZKMTJPQ2NFa2dxLzBMZ2MxYjVmcXJ5ck1KZGU5R3gvNnJzaDFCYmZzVlZvSm1ZSEJZaWZScmNBcDlDR28rcFRLZlU1ajhwQlp2b1Q1QXFBdVdtTlozbUVHcmVWWHIwR1pZQzVnS1VDb0J1ZTNnNkphL0hkajIwODhuWTA3amkyNUd4UVFtcFhINHNFNnZWQTJidXh3T3hPR0FuSTRVM1REWFlweC9qYjUzcUJ6V0FVUmxpR3diYURScVBLK1ZvbGVhRnFPbXJxKzNBQjNJNnJVUWtLN0hVdVVyZzlyUXEwSmxVZ1RoZzZMZ244dEMxbVVmV1pOKzJaZGRrTXBrVUhsaVRxTkd5UnN5d25ZYjB6QVN5WFA2cyswRXJ2MDI3RHpUT1BJcmN4aVVESUZVUGNqcG8zV3BwNTF5YkVheDZPNzVjMVNsa3RhR2hoNGprZUwxeW5aZUs2MGZZU3JybU1DSS85RjJiNU0rY1MwVTZ6Y0dwenRFT2FpWVZTZmtrMXQ2Y0UvMG9aQWxIcmJtS3V5UCswYnBwbEdXMVZWamNVQldOOFBrdmk1ODdRODVvcDVhbTJTRGRxdGNySVNTVTJMdytPaTA5T1M3QmNIQlVYc2ZFUkwyYnQ4Qm5CbGU3bHdmQzNGTEt6RklDNUc1ajBnd281Qlh3NVR5RUVPZnpPUHFXNXd6WSs4dEsvUFdwUVg3dTdmMnBmZVc4WWtlblIxOTZScEFoK1JhcTZoSmpYRDdteUlwdVAwaElUNHg5MzJpd25jSHI2M05TOXk2ZVEvNDRsdmdvaDd3RlBxT1VYYnh5YWNkWE1wa0MvbEQ5MFBNV0luWk9kaUVDSnRWTDgxMHpJNTNZVy8xYjZjSFhycGVESTcwcW9ibEFwTzB3eVNURENuVURxWFJveFFNdElmRlh2YzdFWEtMdU04M1lsTjRoRzhKR2YvZGYzWnQyTGdEbU1uK2xrdy9QWEZ2SytrTXdzaEtJYVJrVTNNekh0SXlhRkRlRTRRQnQxVTJGbEtneWx5QitvSGM5QWkxUHRJTWx4dmNISk9yWGo5VWJSakthZTlOS1lLdUVrcUNydDIvYy9QUzdRdC9uRG9XbUpZY3l5WmVPNzExWThDV2pZQ1hTNVE4SlBMekUva0ZlQ25uSVlOTmV3S3pPZndHNGZQMlZvMWNwTzlSREJrMW95OWZqTnRlVExyNnA5MkcyY20rVVVUcFFBZG5QYWhORjB0S3Zab1FHWjRhZGpkM1J6b0JwQk5BUk1pLzhuTE9zSEJCZXdEWUJnQlkwc0FyT25qRjBMeWtabkdiS3NQcU9DUSsvNEZjU2tlN1dXWWQxMmJrT2l5d3l5NlljTXJldUdXelUvSzVLZGhwaEY4N0JOTnVzZnNWT3RGdm5lbWNteXg4NDRpR29jQ3NCOEQvTk1BbjdxNjhlY0lMZ08xZkFiQ29hdmlvYlY0MmNKZlFtbkphMFlsc1d1eXoraGpPMDN1ODFuczhoQ1FYRmh0VU5MT21ZTEEvWTlCVTZUQS9HemMydVcwOGoxWG1NZlM2ZEMvNmVWTk9HZWJ1SERNbkk1emZvcFBYN0RvR3JvUjdVeTRjM1FEQUx3Q0F0K0thQlZYRFJ4MTNFVzBvcHFWdEo4WWZaMmI2MDFOd1VGWkVZMmx3UzBVU3ltVlpCZHhSdVhoYzNUOXJNTS8wajAvMXpFOElWenh0MkFKdmhFOWhFVTZtWjI2SVR3VkJsOERCbytDWG5lQjY2RDdTdVNQL0JtQXRBR0FLcVppVDFpeW9HMmRrTlpMNmJIb2xrUXNST3JzS3BTcUdVcHhtTTllYnRmV2VJWGhtbUlOOTFNNTJzejZwV1o4c3o5NjcrQ3VUdkhkRzJ2d0lUOGtsS25ucHBZUUxiY1ZKeWticXM4d3IwbEo4MHBuRDN3THdUd0NBQTNvMEFkUGZ5ZGlqZ2dvelhPeHF6NTl2VDM2RE11Zk56eGJWbVN1T1pxT003dEpEejJVWkhTbmg4bCsvbXZjQ2swZlhMYmJrZjVpV09LVTVjL2FPYVUzdDIxYzhKWlEyMWsyZk0zQ2tOTHlXVFlvN2VuQVZBRjhEQUFhS0tiYWFCNTRXbW8xYjRoSXgzcURzT1FOM2Niamo0NGhvZVVLeE10bUR6ZlppMkxDcFBvZDRQb0ViZTF3ZStJMkduZTlrVnkycnFyRGxWOWlNZW5taSs0TlQrSGxjdkdDQjUvb2JYWXBhVjBkSjNQRTk0UDlFNlFTMGdHVElqeDVpWjArSmF4YTBuQVU5ZDhHTVlLL1YyQnNkTnFQRjV2U2Y1L1NTcExDaVpHckpvZTg5MTlZdlNsbDJWUGZwS1JWYnRtR3ovZGhNTCtaV0xUc2tuMTUyZnJMd25PSnFVeTJsT3ZiVTcydC85TiswQVVCeE41RDdjZlVGcVZRYWhVaStRYWJjeU1uRlNUdkxSSDFOVENUekNTKzlXWkpGSzhFMWZ3bm1Ob0Q1SDhEbm5lRHorZTh4dkQ5R1BvTHdpUFZpQ3NRalY4UDNHSTNFb3FyYlZGbzBJVDA0TWk0Z0llbEtOZUZHQnlrYUZCQ1NJKzdjOFVzbTdrbUk4VDU3NEdSa1VIWHRRNllZaW1ZUkk4cWlveXJqY0V4OFZPeHg2TVMyOGJCTlM5RmVTL2dEaThuK1MrVzREL1Q0eXJJWVFrUEt2VG9pc1k2WVZIM24xdU9vOE55SWM4VGdRN2NPSDdwMUxpUW1JdmI2UmVDei83ZW8rSWppSXVwZFlneVRYUUkxTU9qc1BDb3R1ZXBwUmdPU0R3c0wrWEphVjArWlNGdlZxV1VLTk9XQ1BycWd0N3h6b0Vxb1lZcFVwWHhwSWN6UGJvQXo1QkNKVFl2UG90NHNTTGxLVGJ4OEdYZjJSTERmUHQrdDRCOC9mRVBPVHBPcG5tY1VaWFA1TUN4R2FtQjJTbkZTVVYxR1JYTXUxRkhBbFphMEtXZ0NkWVc0dDFLc3FaQm9xNlI5TExHRzFhbG1kRHlud2RJaVRrZGVOU2VUeVNSUUMrT0pHYmNTQ1ZkdVJwOC9HeHI0eDJFdkg1L053SHZMdWtDZjdlSG5BaEp4RjRseFlaVEVxNHlzQk9qQm5lYXkrengycHBSVGdQSnBHaUZkTDYweWRkY01ka09Ec3VyQjdocVRqSzNsbFE0STZib3VWaCtmcWVMUzJ1aHAxWmx4cFlUSW9ydlhzdVBDOEJjRExoemNmR1RuVC84REZPSzlrVllJTytVQUFBQUFTVVZPUks1Q1lJST0ifSwic2x1ZyI6Im15LWZpcnN0LWNvbnRlbnQiLCJjb2xsZWN0aWJsZVdhbGxldEFkZHJlc3MiOiIiLCJoaWdobGlnaHRzQ2hhaW4iOiJvcHRpbWlzbSIsImNvbGxlY3RpYmxlc0Rpc2FibGVkIjp0cnVlLCJwdWJsaXNoZWRBdCI6MTY5NzA5MTM3NjgxNiwicHVibGlzaGVkIjp0cnVlLCJ1cGRhdGVkQXQiOjE2OTcwOTEzNzU2MTIsInBhcmVudElkIjoibUt3RHpXRUhRRG9mOEhwdzlweXIiLCJtYXJrZG93biI6InRoYXQncyBteSBmaXJzdCBwb3N0IGhlcmVcbj09PT09PT09PT09PT09PT09PT09PT09PT1cblxuIVtdKGh0dHBzOi8vc3RvcmFnZS5nb29nbGVhcGlzLmNvbS9wYXB5cnVzX2ltYWdlcy85N2YzN2JlNzQyMjUyYTJkYTUwYWI5YWMwZjNhNDg1MS5qcGcpXG5cbnRoYXQncyB0aGUgY29udGVudCJ9",
						// Tag [{Name:QXBwTmFtZQ Value:UGFyYWdyYXBo} {Name:Q29udGVudC1UeXBl Value:YXBwbGljYXRpb24vanNvbg} {Name:Q29udHJpYnV0b3I Value:MHg1NDJFNEMzYjRhMURDRTBBMUVjYTdCYkMxNDc1NEE4NjdkNjE4NzhB} {Name:UG9zdElk Value:WWVDaWRDVzBNajlRcmo2aEo4Z1g} {Name:Q2F0ZWdvcnk Value:ZGF0YQ} {Name:UG9zdFNsdWc Value:bXktZmlyc3QtY29udGVudA} {Name:UHVibGljYXRpb25JZA Value:ZHFVUTljSnMxc29NQVY3R0lCQjQ} {Name:UHVibGljYXRpb25TbHVn Value:QHl1LXRlc3Q}]
						Tags: []arweave.Tag{
							{Name: "QXBwTmFtZQ", Value: "UGFyYWdyYXBo"},
							{Name: "Q29udGVudC1UeXBl", Value: "YXBwbGljYXRpb24vanNvbg"},
							{Name: "Q29udHJpYnV0b3I", Value: "MHg1NDJFNEMzYjRhMURDRTBBMUVjYTdCYkMxNDc1NEE4NjdkNjE4NzhB"},
							{Name: "UG9zdElk", Value: "WWVDaWRDVzBNajlRcmo2aEo4Z1g"},
							{Name: "Q2F0ZWdvcnk", Value: "ZGF0YQ"},
							{Name: "UG9zdFNsdWc", Value: "bXktZmlyc3QtY29udGVudA"},
							{Name: "UHVibGljYXRpb25JZA", Value: "ZHFVUTljSnMxc29NQVY3R0lCQjQ"},
							{Name: "UHVibGljYXRpb25TbHVn", Value: "QHl1LXRlc3Q"},
						},
					},
				},
				config: &config.Module{
					IPFSGateways: []string{"https://ipfs.rss3.page/"},
				},
			},
			want: &schema.Feed{
				ID:       "Sz5fY8Loj67fWxLQv98r5U5-h2aIA5x4FMsAVP1N2ig",
				Network:  filter.NetworkArweave,
				Index:    0,
				From:     "w5AtiFsNvORfcRtikbdrp2tzqixb05vdPw-ZhgVkD70",
				To:       "w5AtiFsNvORfcRtikbdrp2tzqixb05vdPw-ZhgVkD70",
				Type:     filter.TypeSocialPost,
				Platform: lo.ToPtr(filter.PlatformParagraph),
				Fee: &schema.Fee{
					Amount:  decimal.NewFromInt(212017846),
					Decimal: 12,
				},
				Actions: []*schema.Action{
					{
						Type:     filter.TypeSocialPost,
						Tag:      filter.TagSocial,
						Platform: filter.PlatformParagraph.String(),
						From:     "0x542E4C3b4a1DCE0A1Eca7BbC14754A867d61878A",
						To:       "w5AtiFsNvORfcRtikbdrp2tzqixb05vdPw-ZhgVkD70",
						Metadata: &metadata.SocialPost{
							Handle:  "yu-test",
							Title:   "Yu Test Post",
							Summary: "that's my first post herethat's the content",
							Body:    "that's my first post here\n=========================\n\n![](https://storage.googleapis.com/papyrus_images/97f37be742252a2da50ab9ac0f3a4851.jpg)\n\nthat's the content",
							Media: []metadata.Media{
								{
									Address:  "https://storage.googleapis.com/papyrus_images/f9d95e6eded4d00a4d752f4d004c1c29.jpg",
									MimeType: "image/jpeg",
								},
							},
							ProfileID:     "mKwDzWEHQDof8Hpw9pyr",
							PublicationID: "my-first-content",
							ContentURI:    "https://arweave.net/Sz5fY8Loj67fWxLQv98r5U5-h2aIA5x4FMsAVP1N2ig",
							Tags:          []string{"data"},
							Timestamp:     1697091375,
						},
					},
				},
				Status:    true,
				Timestamp: 1697091466,
			},
			wantError: require.NoError,
		},
		{
			name: "Paragraph Revise",
			arguments: arguments{
				task: &source.Task{
					Network: filter.NetworkArweave,
					Block: arweave.Block{
						Timestamp: 1697092032,
					},
					Transaction: arweave.Transaction{
						ID:       "Xf7C--gk4hlH3mG0UnFiISYgOdymfInv2EgeOF0GeNg",
						Reward:   "212017846",
						Quantity: "0",
						Owner:    "rsfrA1_2H7Pb4kRtHj6EryEELG1sksd-1xGbAWJJqgCIJs9dQYL2C7afuCFX-pryKFpU3ZLssERyObt-BiDwWA3vSHAFljt0CbCBZRKqWKWeEXXdoBLR_Vf8724P14YqRubW7a0n6UaZKsJsxah35yPCANnw9QbnHJouTlNyky41ZnbBClRlYWr1_PkEMvFEsQcqIE5J8jcgJlaTNtiOi7ruvRP3z-NtqufuJFFq3_4hrL6ICpbJnZBgZuX33tr6YvCrYExtFmd8wJoL4s6MSioKYSWYk60ngr8EgUHotS1lzPemWRhY9wjbrg3wh00sCO44wv5CmE2Ke-EoKZYKrUu8g5z2MlPwOnNxBj67wzrSRzkpaVvbEWqneEwG_UcDxKU_SLeJ0_qGLNkQjgqjhfAAEdivsfV0Fz3hNRmVu2ae84QtoPQyvvcr2JLe-bTjbGvna_C52fR7-p9sp-MlZnL8vPnKfPZrTvfCOd935O2_CdiyzvOA35jQKQhe5UhqwH0hoYdplE2DHRN6MR42n-8nq3vqxp7Y34l-aUxnRIHBquMFbfH4KKn8N322_e_6nAwImjp_DziPhz5xOyQJgZOzCBTFuQrbaHkGbQ6ou814fyAUDJlA3S5-WKtsD8Jk1AMg0YmIdFUgCVUwwepoAgK1UPAxpq64GouKmnqjI58",
						Data:     "eyJzdG9yZU9uQXJ3ZWF2ZSI6dHJ1ZSwidGl0bGUiOiJZdSBUZXN0IFBvc3QiLCJhcndlYXZlSWQiOiJTejVmWThMb2o2N2ZXeExRdjk4cjVVNS1oMmFJQTV4NEZNc0FWUDFOMmlnIiwiaGlnaGxpZ2h0c0NoYWluIjoib3B0aW1pc20iLCJ1Z2MiOnRydWUsImNvbGxlY3RpYmxlc0Rpc2FibGVkIjp0cnVlLCJhcmNoaXZlZCI6ZmFsc2UsImNyZWF0ZWRBdCI6MTY5NzA5MTYwMDU3NywicG9zdF9wcmV2aWV3IjoidGhhdCdzIG15IGZpcnN0IHBvc3QgaGVyZXRoYXQncyB0aGUgY29udGVudCIsImlkIjoiZzZHMmxoRG1TbE5OQ3BlcGltREUiLCJjYXRlZ29yaWVzIjpbImRhdGEiXSwic2x1ZyI6Im15LWZpcnN0LWNvbnRlbnQiLCJzZW5kWE1UUCI6ZmFsc2UsInB1Ymxpc2hlZEF0IjoxNjk3MDkxMzc2ODE2LCJzZW5kTmV3c2xldHRlciI6ZmFsc2UsImRvbnRQdWJsaXNoT25saW5lIjpmYWxzZSwicGFyZW50SWQiOiJtS3dEeldFSFFEb2Y4SHB3OXB5ciIsImFjY2Vzc1Jlc3RyaWN0aW9uIjpudWxsLCJzdWJ0aXRsZSI6IlBvc3QiLCJjb3Zlcl9pbWciOnsiaW1nIjp7InNyYyI6Imh0dHBzOi8vc3RvcmFnZS5nb29nbGVhcGlzLmNvbS9wYXB5cnVzX2ltYWdlcy9mOWQ5NWU2ZWRlZDRkMDBhNGQ3NTJmNGQwMDRjMWMyOS5qcGciLCJ3aWR0aCI6MjA0OCwiaGVpZ2h0IjoxMzU3fSwiaXNIZXJvIjp0cnVlLCJiYXNlNjQiOiJkYXRhOmltYWdlL3BuZztiYXNlNjQsaVZCT1J3MEtHZ29BQUFBTlNVaEVVZ0FBQUNBQUFBQVZDQUlBQUFDb3IzdTlBQUFBQ1hCSVdYTUFBQXNUQUFBTEV3RUFtcHdZQUFBSDYwbEVRVlI0bkNYQmVWVFNlUUlBOEcrenpiemROelB2Ylc5MzUvaGpwM2JxZFV4VFRiNXlLaXQxSisyd083TkRMY3NvNzlBU3NCUVZ4YlBTRkFVRlJmMVpHb28vL1lrZ2Nrb2d2MFFVa0t1VVE4VVFqN1FtdFZMNzdSLzcrWUN2UTJrQlZ6S1RMbDRtWDc1MHBoeFpyVmo0RXYxenRXSm10V0x5TC9LeE5STGJWNktYNjRYVy9TS2psOGpvS3pINUM3WDdCWnJEUXMwQmtkNUhhTnpOMTIzaDZmNk9ESzVDektESkFCb05nR01GMVdidytBWElVWUJVSWZqZDIrZTcwTHp0Y1JWSjEyT0NjRWtncS8wTGdjMWI1ZnFyeXJNSmRlOUd4LzZyc2gxQmJmc1ZWb0ptWUhCWWlmUnJjQXA5Q0dvK3BUS2ZVNWo4cEJadm9UNUFxQXVXbU5aM21FR3JlVlhyMEdaWUM1Z0tVQ29CdWUzZzZKYS9IZGoyMDg4blkwN2ppMjVHeFFRbXBYSDRzRTZ2VkEyYnV4d094T0dBbkk0VTNURFhZcHgvamI1M3FCeldBVVJsaUd3YmFEUnFQSytWb2xlYUZxT21ycSszQUIzSTZyVVFrSzdIVXVVcmc5clFxMEpsVWdUaGc2TGduOHRDMW1VZldaTisyWmRka01wa1VIbGlUcU5HeVJzeXduWWIwekFTeVhQNnMrMEVydjAyN0R6VE9QSXJjeGlVRElGVVBjanBvM1dwcDUxeWJFYXg2Tzc1YzFTbGt0YUdoaDRqa2VMMXluWmVLNjBmWVNycm1NQ0kvOUYyYjVNK2NTMFU2emNHcHp0RU9haVlWU2ZrazF0NmNFLzBvWkFsSHJibUt1eVArMGJwcGxHVzFWVmpjVUJXTjhQa3ZpNTg3UTg1b3A1YW0yU0RkcXRjcklTU1UyTHcrT2kwOU9TN0JjSEJVWHNmRVJMMmJ0OEJuQmxlN2x3ZkMzRkxLekZJQzVHNWowZ3dvNUJYdzVUeUVFT2Z6T1BxVzV3elkrOHRLL1BXcFFYN3U3ZjJwZmVXOFlrZW5SMTk2UnBBaCtSYXE2aEpqWEQ3bXlJcHVQMGhJVDR4OTMyaXduY0hyNjNOUzl5NmVRLzQ0bHZnb2g3d0ZQcU9VWGJ4eWFjZFhNcGtDL2xEOTBQTVdJblpPZGlFQ0p0Vkw4MTB6STUzWVcvMWI2Y0hYcnBlREk3MHFvYmxBcE8wd3lTVERDblVEcVhSb3hRTXRJZkZYdmM3RVhLTHVNODNZbE40aEc4SkdmL2RmM1p0MkxnRG1Nbitsa3cvUFhGdksra013c2hLSWFSa1UzTXpIdEl5YUZEZUU0UUJ0MVUyRmxLZ3lseUIrb0hjOUFpMVB0SU1seHZjSEpPclhqOVViUmpLYWU5TktZS3VFa3FDcnQyL2MvUFM3UXQvbkRvV21KWWN5eVplTzcxMVk4Q1dqWUNYUzVROEpQTHpFL2tGZUNubklZTk5ld0t6T2Z3RzRmUDJWbzFjcE85UkRCazFveTlmak50ZVRMcjZwOTJHMmNtK1VVVHBRQWRuUGFoTkYwdEt2Wm9RR1o0YWRqZDNSem9CcEJOQVJNaS84bkxPc0hCQmV3RFlCZ0JZMHNBck9uakYwTHlrWm5HYktzUHFPQ1ErLzRGY1NrZTdXV1lkMTJia09peXd5eTZZY01yZXVHV3pVL0s1S2RocGhGODdCTk51c2ZzVk90RnZuZW1jbXl4ODQ0aUdvY0NzQjhEL05NQW43cTY4ZWNJTGdPMWZBYkNvYXZpb2JWNDJjSmZRbW5KYTBZbHNXdXl6K2hqTzAzdTgxbnM4aENRWEZodFVOTE9tWUxBL1k5QlU2VEEvR3pjMnVXMDhqMVhtTWZTNmRDLzZlVk5PR2VidUhETW5JNXpmb3BQWDdEb0dyb1I3VXk0YzNRREFMd0NBdCtLYUJWWERSeDEzRVcwb3BxVnRKOFlmWjJiNjAxTndVRlpFWTJsd1MwVVN5bVZaQmR4UnVYaGMzVDlyTU0vMGowLzF6RThJVnp4dDJBSnZoRTloRVU2bVoyNklUd1ZCbDhEQm8rQ1huZUI2NkQ3U3VTUC9CbUF0QUdBS3FaaVQxaXlvRzJka05aTDZiSG9sa1FzUk9yc0twU3FHVXB4bU05ZWJ0ZldlSVhobW1JTjkxTTUyc3o2cFdaOHN6OTY3K0N1VHZIZEcydndJVDhrbEtubnBwWVFMYmNWSnlrYnFzOHdyMGxKODBwbkQzd0x3VHdDQUEzbzBBZFBmeWRpamdnb3pYT3hxejU5dlQzNkRNdWZOenhiVm1TdU9acU9NN3RKRHoyVVpIU25oOGwrL212Y0NrMGZYTGJia2Y1aVdPS1U1Yy9hT2FVM3QyMWM4SlpRMjFrMmZNM0NrTkx5V1RZbzdlbkFWQUY4REFBYUtLYmFhQjU0V21vMWI0aEl4M3FEc09RTjNjYmpqNDRob2VVS3hNdG1EemZaaTJMQ3BQb2Q0UG9FYmUxd2UrSTJHbmU5a1Z5MnJxckRsVjlpTWVubWkrNE5UK0hsY3ZHQ0I1L29iWFlwYVYwZEozUEU5NFA5RTZRUzBnR1RJang1aVowK0pheGEwbkFVOWQ4R01ZSy9WMkJzZE5xUEY1dlNmNS9TU3BMQ2laR3JKb2U4OTE5WXZTbGwyVlBmcEtSVmJ0bUd6L2RoTUwrWldMVHNrbjE1MmZyTHduT0pxVXkybE92YlU3MnQvOU4rMEFVQnhONUQ3Y2ZVRnFWUWFoVWkrUWFiY3lNbkZTVHZMUkgxTlRDVHpDUys5V1pKRks4RTFmd25tTm9ENUg4RG5uZUR6K2U4eHZEOUdQb0x3aVBWaUNzUWpWOFAzR0kzRW9xcmJWRm8wSVQwNE1pNGdJZWxLTmVGR0J5a2FGQkNTSSs3YzhVc203a21JOFQ1NzRHUmtVSFh0UTZZWWltWVJJOHFpb3lyamNFeDhWT3h4Nk1TMjhiQk5TOUZlUy9nRGk4bitTK1c0RC9UNHlySVlRa1BLdlRvaXNZNllWSDNuMXVPbzhOeUljOFRnUTdjT0g3cDFMaVFtSXZiNlJlQ3ovN2VvK0lqaUl1cGRZZ3lUWFFJMU1PanNQQ290dWVwcFJnT1NEd3NMK1hKYVYwK1pTRnZWcVdVS05PV0NQcnFndDd4em9FcW9ZWXBVcFh4cEljelBib0F6NUJDSlRZdlBvdDRzU0xsS1RieDhHWGYyUkxEZlB0K3Q0QjgvZkVQT1RwT3BubWNVWlhQNU1DeEdhbUIyU25GU1VWMUdSWE11MUZIQWxaYTBLV2dDZFlXNHQxS3NxWkJvcTZSOUxMR0cxYWxtZER5bndkSWlUa2RlTlNlVHlTUlFDK09KR2JjU0NWZHVScDgvR3hyNHgyRXZINS9Od0h2THVrQ2Y3ZUhuQWhKeEY0bHhZWlRFcTR5c0JPakJuZWF5K3p4MnBwUlRnUEpwR2lGZEw2MHlkZGNNZGtPRHN1ckI3aHFUakszbGxRNEk2Ym91VmgrZnFlTFMydWhwMVpseHBZVElvcnZYc3VQQzhCY0RMaHpjZkdUblQvOERGT0s5a1ZZSU8rVUFBQUFBU1VWT1JLNUNZSUk9In0sImF1dGhvcnMiOlsibUt3RHpXRUhRRG9mOEhwdzlweXIiXSwiY29sbGVjdGlibGVXYWxsZXRBZGRyZXNzIjoiIiwianNvbiI6IntcInR5cGVcIjpcImRvY1wiLFwiY29udGVudFwiOlt7XCJ0eXBlXCI6XCJoZWFkaW5nXCIsXCJhdHRyc1wiOntcInRleHRBbGlnblwiOlwibGVmdFwiLFwibGV2ZWxcIjoxfSxcImNvbnRlbnRcIjpbe1widHlwZVwiOlwidGV4dFwiLFwidGV4dFwiOlwidGhhdCdzIG15IHNlY29uZCBwb3N0IGhlcmUocmV2aXNlZClcIn1dfSx7XCJ0eXBlXCI6XCJmaWd1cmVcIixcImF0dHJzXCI6e1wic3JjXCI6bnVsbCxcImZpbGVcIjpudWxsLFwiYWx0XCI6bnVsbCxcInRpdGxlXCI6bnVsbCxcImJsdXJEYXRhVVJMXCI6bnVsbCxcImZsb2F0XCI6XCJub25lXCIsXCJ3aWR0aFwiOm51bGx9LFwiY29udGVudFwiOlt7XCJ0eXBlXCI6XCJpbWFnZVwiLFwiYXR0cnNcIjp7XCJzcmNcIjpcImh0dHBzOi8vc3RvcmFnZS5nb29nbGVhcGlzLmNvbS9wYXB5cnVzX2ltYWdlcy85N2YzN2JlNzQyMjUyYTJkYTUwYWI5YWMwZjNhNDg1MS5qcGdcIixcImFsdFwiOm51bGwsXCJ0aXRsZVwiOm51bGwsXCJibHVyZGF0YXVybFwiOlwiZGF0YTppbWFnZS9wbmc7YmFzZTY0LGlWQk9SdzBLR2dvQUFBQU5TVWhFVWdBQUFDQUFBQUFWQ0FJQUFBQ29yM3U5QUFBQUNYQklXWE1BQUFzVEFBQUxFd0VBbXB3WUFBQUg2MGxFUVZSNG5DWEJlVlRTZVFJQThHK3p6YnpkTnpQdmJXOTM1L2hqcDNicWRVeFRUYjV5S2l0MUorMndPN05ETGNzbzc5QVNzQlFWeGJQU0ZBVUZSZjFaR29vLy9Za2dja29ndjBRVWtLdVVROFVRajdRbXRWTDc3Ui83K1lDdlEya0JWektUTGw0bVg3NTBwaHhaclZqNEV2MXp0V0ptdFdMeUwvS3hOUkxiVjZLWDY0WFcvU0tqbDhqb0t6SDVDN1g3QlpyRFFzMEJrZDVIYU56TjEyM2g2ZjZPREs1Q3pLREpBQm9OZ0dNRjFXYncrQVhJVVlCVUlmamQyK2U3MEx6dGNSVkoxMk9DY0VrZ3EvMExnYzFiNWZxcnlyTUpkZTlHeC82cnNoMUJiZnNWVm9KbVlIQllpZlJyY0FwOUNHbytwVEtmVTVqOHBCWnZvVDVBcUF1V21OWjNtRUdyZVZYcjBHWllDNWdLVUNvQnVlM2c2SmEvSGRqMjA4OG5ZMDdqaTI1R3hRUW1wWEg0c0U2dlZBMmJ1eHdPeE9HQW5JNFUzVERYWXB4L2piNTNxQnpXQVVSbGlHd2JhRFJxUEsrVm9sZWFGcU9tcnErM0FCM0k2clVRa0s3SFV1VXJnOXJRcTBKbFVnVGhnNkxnbjh0QzFtVWZXWk4rMlpkZGtNcGtVSGxpVHFOR3lSc3l3blliMHpBU3lYUDZzKzBFcnYwMjdEelRPUElyY3hpVURJRlVQY2pwbzNXcHA1MXliRWF4Nk83NWMxU2xrdGFHaGg0amtlTDF5blplSzYwZllTcnJtTUNJLzlGMmI1TStjUzBVNnpjR3B6dEVPYWlZVlNma2sxdDZjRS8wb1pBbEhyYm1LdXlQKzBicHBsR1cxVlZqY1VCV044UGt2aTU4N1E4NW9wNWFtMlNEZHF0Y3JJU1NVMkx3K09pMDlPUzdCY0hCVVhzZkVSTDJidDhCbkJsZTdsd2ZDM0ZMS3pGSUM1RzVqMGd3bzVCWHc1VHlFRU9mek9QcVc1d3pZKzh0Sy9QV3BRWDd1N2YycGZlVzhZa2VuUjE5NlJwQWgrUmFxNmhKalhEN215SXB1UDBoSVQ0eDkzMml3bmNIcjYzTlM5eTZlUS80NGx2Z29oN3dGUHFPVVhieHlhY2RYTXBrQy9sRDkwUE1XSW5aT2RpRUNKdFZMODEwekk1M1lXLzFiNmNIWHJwZURJNzBxb2JsQXBPMHd5U1REQ25VRHFYUm94UU10SWZGWHZjN0VYS0x1TTgzWWxONGhHOEpHZi9kZjNadDJMZ0RtTW4rbGt3L1BYRnZLK2tNd3NoS0lhUmtVM016SHRJeWFGRGVFNFFCdDFVMkZsS2d5bHlCK29IYzlBaTFQdElNbHh2Y0hKT3JYajlVYlJqS2FlOU5LWUt1RWtxQ3J0Mi9jL1BTN1F0L25Eb1dtSlljeXlaZU83MTFZOENXallDWFM1UThKUEx6RS9rRmVDbm5JWU5OZXdLek9md0c0ZlAyVm8xY3BPOVJEQmsxb3k5ZmpOdGVUTHI2cDkyRzJjbStVVVRwUUFkblBhaE5GMHRLdlpvUUdaNGFkamQzUnpvQnBCTkFSTWkvOG5MT3NIQkJld0RZQmdCWTBzQXJPbmpGMEx5a1puR2JLc1BxT0NRKy80RmNTa2U3V1dZZDEyYmtPaXl3eXk2WWNNcmV1R1d6VS9LNUtkaHBoRjg3Qk5OdXNmc1ZPdEZ2bmVtY215eDg0NGlHb2NDc0I4RC9OTUFuN3E2OGVjSUxnTzFmQWJDb2F2aW9iVjQyY0pmUW1uSmEwWWxzV3V5eitoak8wM3U4MW5zOGhDUVhGaHRVTkxPbVlMQS9ZOUJVNlRBL0d6YzJ1VzA4ajFYbU1mUzZkQy82ZVZOT0dlYnVIRE1uSTV6Zm9wUFg3RG9Hcm9SN1V5NGMzUURBTHdDQXQrS2FCVlhEUngxM0VXMG9wcVZ0SjhZZloyYjYwMU53VUZaRVkybHdTMFVTeW1WWkJkeFJ1WGhjM1Q5ck1NLzBqMC8xekU4SVZ6eHQyQUp2aEU5aEVVNm1aMjZJVHdWQmw4REJvK0NYbmVCNjZEN1N1U1AvQm1BdEFHQUtxWmlUMWl5b0cyZGtOWkw2YkhvbGtRc1JPcnNLcFNxR1VweG1NOWVidGZXZUlYaG1tSU45MU01MnN6NnBXWjhzejk2NytDdVR2SGRHMnZ3SVQ4a2xLbm5wcFlRTGJjVkp5a2Jxczh3cjBsSjgwcG5EM3dMd1R3Q0FBM28wQWRQZnlkaWpnZ296WE94cXo1OXZUMzZETXVmTnp4YlZtU3VPWnFPTTd0SkR6MlVaSFNuaDhsKy9tdmNDazBmWExiYmtmNWlXT0tVNWMvYU9hVTN0MjFjOEpaUTIxazJmTTNDa05MeVdUWW83ZW5BVkFGOERBQWFLS2JhYUI1NFdtbzFiNGhJeDNxRHNPUU4zY2JqajQ0aG9lVUt4TXRtRHpmWmkyTENwUG9kNFBvRWJlMXdlK0kyR25lOWtWeTJycXJEbFY5aU1lbm1pKzROVCtIbGN2R0NCNS9vYlhZcGFWMGRKM1BFOTRQOUU2UVMwZ0dUSWp4NWlaMCtKYXhhMG5BVTlkOEdNWUsvVjJCc2ROcVBGNXZTZjUvU1NwTENpWkdySm9lODkxOVl2U2xsMlZQZnBLUlZidG1Hei9kaE1MK1pXTFRza24xNTJmckx3bk9KcVV5MmxPdmJVNzJ0LzlOKzBBVUJ4TjVEN2NmVUZxVlFhaFVpK1FhYmN5TW5GU1R2TFJIMU5UQ1R6Q1MrOVdaSkZLOEUxZndubU5vRDVIOERubmVEeitlOHh2RDlHUG9Md2lQVmlDc1FqVjhQM0dJM0VvcXJiVkZvMElUMDRNaTRnSWVsS05lRkdCeWthRkJDU0krN2M4VXNtN2ttSThUNTc0R1JrVUhYdFE2WVlpbVlSSThxaW95cmpjRXg4Vk94eDZNUzI4YkJOUzlGZVMvZ0RpOG4rUytXNEQvVDR5cklZUWtQS3ZUb2lzWTZZVkgzbjF1T284TnlJYzhUZ1E3Y09IN3AxTGlRbUl2YjZSZUN6LzdlbytJamlJdXBkWWd5VFhRSTFNT2pzUENvdHVlcHBSZ09TRHdzTCtYSmFWMCtaU0Z2VnFXVUtOT1dDUHJxZ3Q3eHpvRXFvWVlwVXBYeHBJY3pQYm9BejVCQ0pUWXZQb3Q0c1NMbEtUYng4R1hmMlJMRGZQdCt0NEI4L2ZFUE9UcE9wbm1jVVpYUDVNQ3hHYW1CMlNuRlNVVjFHUlhNdTFGSEFsWmEwS1dnQ2RZVzR0MUtzcVpCb3E2UjlMTEdHMWFsbWREeW53ZElpVGtkZU5TZVR5U1JRQytPSkdiY1NDVmR1UnA4L0d4cjR4MkV2SDUvTndIdkx1a0NmN2VIbkFoSnhGNGx4WVpURXE0eXNCT2pCbmVheSt6eDJwcFJUZ1BKcEdpRmRMNjB5ZGRjTWRrT0RzdXJCN2hxVGpLM2xsUTRJNmJvdVZoK2ZxZUxTMnVocDFabHhwWVRJb3J2WHN1UEM4QmNETGh6Y2ZHVG5ULzhERk9LOWtWWUlPK1VBQUFBQVNVVk9SSzVDWUlJPVwiLFwibmV4dGhlaWdodFwiOjEzNTcsXCJuZXh0d2lkdGhcIjoyMDQ4fX0se1widHlwZVwiOlwiZmlnY2FwdGlvblwifV19LHtcInR5cGVcIjpcInBhcmFncmFwaFwiLFwiYXR0cnNcIjp7XCJ0ZXh0QWxpZ25cIjpcImxlZnRcIn0sXCJjb250ZW50XCI6W3tcInR5cGVcIjpcInRleHRcIixcInRleHRcIjpcInRoYXQncyB0aGUgY29udGVudFwifV19XX0iLCJzdGF0aWNIdG1sIjoiPGgxPnRoYXQncyBteSBzZWNvbmQgcG9zdCBoZXJlKHJldmlzZWQpPC9oMT48ZmlndXJlIGZsb2F0PVwibm9uZVwiIGRhdGEtdHlwZT1cImZpZ3VyZVwiIGNsYXNzPVwiaW1nLWNlbnRlclwiIHN0eWxlPVwibWF4LXdpZHRoOiBudWxsO1wiPjxpbWcgc3JjPVwiaHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL3BhcHlydXNfaW1hZ2VzLzk3ZjM3YmU3NDIyNTJhMmRhNTBhYjlhYzBmM2E0ODUxLmpwZ1wiIGJsdXJkYXRhdXJsPVwiZGF0YTppbWFnZS9wbmc7YmFzZTY0LGlWQk9SdzBLR2dvQUFBQU5TVWhFVWdBQUFDQUFBQUFWQ0FJQUFBQ29yM3U5QUFBQUNYQklXWE1BQUFzVEFBQUxFd0VBbXB3WUFBQUg2MGxFUVZSNG5DWEJlVlRTZVFJQThHK3p6YnpkTnpQdmJXOTM1L2hqcDNicWRVeFRUYjV5S2l0MUorMndPN05ETGNzbzc5QVNzQlFWeGJQU0ZBVUZSZjFaR29vLy9Za2dja29ndjBRVWtLdVVROFVRajdRbXRWTDc3Ui83K1lDdlEya0JWektUTGw0bVg3NTBwaHhaclZqNEV2MXp0V0ptdFdMeUwvS3hOUkxiVjZLWDY0WFcvU0tqbDhqb0t6SDVDN1g3QlpyRFFzMEJrZDVIYU56TjEyM2g2ZjZPREs1Q3pLREpBQm9OZ0dNRjFXYncrQVhJVVlCVUlmamQyK2U3MEx6dGNSVkoxMk9DY0VrZ3EvMExnYzFiNWZxcnlyTUpkZTlHeC82cnNoMUJiZnNWVm9KbVlIQllpZlJyY0FwOUNHbytwVEtmVTVqOHBCWnZvVDVBcUF1V21OWjNtRUdyZVZYcjBHWllDNWdLVUNvQnVlM2c2SmEvSGRqMjA4OG5ZMDdqaTI1R3hRUW1wWEg0c0U2dlZBMmJ1eHdPeE9HQW5JNFUzVERYWXB4L2piNTNxQnpXQVVSbGlHd2JhRFJxUEsrVm9sZWFGcU9tcnErM0FCM0k2clVRa0s3SFV1VXJnOXJRcTBKbFVnVGhnNkxnbjh0QzFtVWZXWk4rMlpkZGtNcGtVSGxpVHFOR3lSc3l3blliMHpBU3lYUDZzKzBFcnYwMjdEelRPUElyY3hpVURJRlVQY2pwbzNXcHA1MXliRWF4Nk83NWMxU2xrdGFHaGg0amtlTDF5blplSzYwZllTcnJtTUNJLzlGMmI1TStjUzBVNnpjR3B6dEVPYWlZVlNma2sxdDZjRS8wb1pBbEhyYm1LdXlQKzBicHBsR1cxVlZqY1VCV044UGt2aTU4N1E4NW9wNWFtMlNEZHF0Y3JJU1NVMkx3K09pMDlPUzdCY0hCVVhzZkVSTDJidDhCbkJsZTdsd2ZDM0ZMS3pGSUM1RzVqMGd3bzVCWHc1VHlFRU9mek9QcVc1d3pZKzh0Sy9QV3BRWDd1N2YycGZlVzhZa2VuUjE5NlJwQWgrUmFxNmhKalhEN215SXB1UDBoSVQ0eDkzMml3bmNIcjYzTlM5eTZlUS80NGx2Z29oN3dGUHFPVVhieHlhY2RYTXBrQy9sRDkwUE1XSW5aT2RpRUNKdFZMODEwekk1M1lXLzFiNmNIWHJwZURJNzBxb2JsQXBPMHd5U1REQ25VRHFYUm94UU10SWZGWHZjN0VYS0x1TTgzWWxONGhHOEpHZi9kZjNadDJMZ0RtTW4rbGt3L1BYRnZLK2tNd3NoS0lhUmtVM016SHRJeWFGRGVFNFFCdDFVMkZsS2d5bHlCK29IYzlBaTFQdElNbHh2Y0hKT3JYajlVYlJqS2FlOU5LWUt1RWtxQ3J0Mi9jL1BTN1F0L25Eb1dtSlljeXlaZU83MTFZOENXallDWFM1UThKUEx6RS9rRmVDbm5JWU5OZXdLek9md0c0ZlAyVm8xY3BPOVJEQmsxb3k5ZmpOdGVUTHI2cDkyRzJjbStVVVRwUUFkblBhaE5GMHRLdlpvUUdaNGFkamQzUnpvQnBCTkFSTWkvOG5MT3NIQkJld0RZQmdCWTBzQXJPbmpGMEx5a1puR2JLc1BxT0NRKy80RmNTa2U3V1dZZDEyYmtPaXl3eXk2WWNNcmV1R1d6VS9LNUtkaHBoRjg3Qk5OdXNmc1ZPdEZ2bmVtY215eDg0NGlHb2NDc0I4RC9OTUFuN3E2OGVjSUxnTzFmQWJDb2F2aW9iVjQyY0pmUW1uSmEwWWxzV3V5eitoak8wM3U4MW5zOGhDUVhGaHRVTkxPbVlMQS9ZOUJVNlRBL0d6YzJ1VzA4ajFYbU1mUzZkQy82ZVZOT0dlYnVIRE1uSTV6Zm9wUFg3RG9Hcm9SN1V5NGMzUURBTHdDQXQrS2FCVlhEUngxM0VXMG9wcVZ0SjhZZloyYjYwMU53VUZaRVkybHdTMFVTeW1WWkJkeFJ1WGhjM1Q5ck1NLzBqMC8xekU4SVZ6eHQyQUp2aEU5aEVVNm1aMjZJVHdWQmw4REJvK0NYbmVCNjZEN1N1U1AvQm1BdEFHQUtxWmlUMWl5b0cyZGtOWkw2YkhvbGtRc1JPcnNLcFNxR1VweG1NOWVidGZXZUlYaG1tSU45MU01MnN6NnBXWjhzejk2NytDdVR2SGRHMnZ3SVQ4a2xLbm5wcFlRTGJjVkp5a2Jxczh3cjBsSjgwcG5EM3dMd1R3Q0FBM28wQWRQZnlkaWpnZ296WE94cXo1OXZUMzZETXVmTnp4YlZtU3VPWnFPTTd0SkR6MlVaSFNuaDhsKy9tdmNDazBmWExiYmtmNWlXT0tVNWMvYU9hVTN0MjFjOEpaUTIxazJmTTNDa05MeVdUWW83ZW5BVkFGOERBQWFLS2JhYUI1NFdtbzFiNGhJeDNxRHNPUU4zY2JqajQ0aG9lVUt4TXRtRHpmWmkyTENwUG9kNFBvRWJlMXdlK0kyR25lOWtWeTJycXJEbFY5aU1lbm1pKzROVCtIbGN2R0NCNS9vYlhZcGFWMGRKM1BFOTRQOUU2UVMwZ0dUSWp4NWlaMCtKYXhhMG5BVTlkOEdNWUsvVjJCc2ROcVBGNXZTZjUvU1NwTENpWkdySm9lODkxOVl2U2xsMlZQZnBLUlZidG1Hei9kaE1MK1pXTFRza24xNTJmckx3bk9KcVV5MmxPdmJVNzJ0LzlOKzBBVUJ4TjVEN2NmVUZxVlFhaFVpK1FhYmN5TW5GU1R2TFJIMU5UQ1R6Q1MrOVdaSkZLOEUxZndubU5vRDVIOERubmVEeitlOHh2RDlHUG9Md2lQVmlDc1FqVjhQM0dJM0VvcXJiVkZvMElUMDRNaTRnSWVsS05lRkdCeWthRkJDU0krN2M4VXNtN2ttSThUNTc0R1JrVUhYdFE2WVlpbVlSSThxaW95cmpjRXg4Vk94eDZNUzI4YkJOUzlGZVMvZ0RpOG4rUytXNEQvVDR5cklZUWtQS3ZUb2lzWTZZVkgzbjF1T284TnlJYzhUZ1E3Y09IN3AxTGlRbUl2YjZSZUN6LzdlbytJamlJdXBkWWd5VFhRSTFNT2pzUENvdHVlcHBSZ09TRHdzTCtYSmFWMCtaU0Z2VnFXVUtOT1dDUHJxZ3Q3eHpvRXFvWVlwVXBYeHBJY3pQYm9BejVCQ0pUWXZQb3Q0c1NMbEtUYng4R1hmMlJMRGZQdCt0NEI4L2ZFUE9UcE9wbm1jVVpYUDVNQ3hHYW1CMlNuRlNVVjFHUlhNdTFGSEFsWmEwS1dnQ2RZVzR0MUtzcVpCb3E2UjlMTEdHMWFsbWREeW53ZElpVGtkZU5TZVR5U1JRQytPSkdiY1NDVmR1UnA4L0d4cjR4MkV2SDUvTndIdkx1a0NmN2VIbkFoSnhGNGx4WVpURXE0eXNCT2pCbmVheSt6eDJwcFJUZ1BKcEdpRmRMNjB5ZGRjTWRrT0RzdXJCN2hxVGpLM2xsUTRJNmJvdVZoK2ZxZUxTMnVocDFabHhwWVRJb3J2WHN1UEM4QmNETGh6Y2ZHVG5ULzhERk9LOWtWWUlPK1VBQUFBQVNVVk9SSzVDWUlJPVwiIG5leHRoZWlnaHQ9XCIxMzU3XCIgbmV4dHdpZHRoPVwiMjA0OFwiIGNsYXNzPVwiaW1hZ2Utbm9kZSBlbWJlZFwiPjxmaWdjYXB0aW9uIGh0bWxhdHRyaWJ1dGVzPVwiW29iamVjdCBPYmplY3RdXCIgY2xhc3M9XCJoaWRlLWZpZ2NhcHRpb25cIj48L2ZpZ2NhcHRpb24-PC9maWd1cmU-PHA-dGhhdCdzIHRoZSBjb250ZW50PC9wPiIsInB1Ymxpc2hlZCI6dHJ1ZSwidXBkYXRlZEF0IjoxNjk3MDkxNjI5NzAzLCJtYXJrZG93biI6InRoYXQncyBteSBzZWNvbmQgcG9zdCBoZXJlKHJldmlzZWQpXG49PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PVxuXG4hW10oaHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL3BhcHlydXNfaW1hZ2VzLzk3ZjM3YmU3NDIyNTJhMmRhNTBhYjlhYzBmM2E0ODUxLmpwZylcblxudGhhdCdzIHRoZSBjb250ZW50In0",
						Tags: []arweave.Tag{
							{Name: "QXBwTmFtZQ", Value: "UGFyYWdyYXBo"},
							{Name: "Q29udGVudC1UeXBl", Value: "YXBwbGljYXRpb24vanNvbg"},
							{Name: "Q29udHJpYnV0b3I", Value: "MHg1NDJFNEMzYjRhMURDRTBBMUVjYTdCYkMxNDc1NEE4NjdkNjE4NzhB"},
							{Name: "UG9zdElk", Value: "ZzZHMmxoRG1TbE5OQ3BlcGltREU"},
							{Name: "Q2F0ZWdvcnk", Value: "ZGF0YQ"},
							{Name: "UG9zdFNsdWc", Value: "bXktZmlyc3QtY29udGVudA"},
							{Name: "UHVibGljYXRpb25JZA", Value: "ZHFVUTljSnMxc29NQVY3R0lCQjQ"},
							{Name: "UHVibGljYXRpb25TbHVn", Value: "QHl1LXRlc3Q"},
						},
					},
				},
				config: &config.Module{
					IPFSGateways: []string{"https://ipfs.rss3.page/"},
				},
			},
			want: &schema.Feed{
				ID:       "Xf7C--gk4hlH3mG0UnFiISYgOdymfInv2EgeOF0GeNg",
				Network:  filter.NetworkArweave,
				Index:    0,
				From:     "w5AtiFsNvORfcRtikbdrp2tzqixb05vdPw-ZhgVkD70",
				To:       "w5AtiFsNvORfcRtikbdrp2tzqixb05vdPw-ZhgVkD70",
				Type:     filter.TypeSocialRevise,
				Platform: lo.ToPtr(filter.PlatformParagraph),
				Fee: &schema.Fee{
					Amount:  decimal.NewFromInt(212017846),
					Decimal: 12,
				},

				Actions: []*schema.Action{
					{
						Type:     filter.TypeSocialRevise,
						Tag:      filter.TagSocial,
						Platform: filter.PlatformParagraph.String(),
						From:     "0x542E4C3b4a1DCE0A1Eca7BbC14754A867d61878A",
						To:       "w5AtiFsNvORfcRtikbdrp2tzqixb05vdPw-ZhgVkD70",
						Metadata: &metadata.SocialPost{
							Handle:  "yu-test",
							Title:   "Yu Test Post",
							Summary: "that's my first post herethat's the content",
							Body:    "that's my second post here(revised)\n===================================\n\n![](https://storage.googleapis.com/papyrus_images/97f37be742252a2da50ab9ac0f3a4851.jpg)\n\nthat's the content",
							Media: []metadata.Media{
								{
									Address:  "https://storage.googleapis.com/papyrus_images/f9d95e6eded4d00a4d752f4d004c1c29.jpg",
									MimeType: "image/jpeg",
								},
							},
							ProfileID:     "mKwDzWEHQDof8Hpw9pyr",
							PublicationID: "my-first-content",
							ContentURI:    "https://arweave.net/Xf7C--gk4hlH3mG0UnFiISYgOdymfInv2EgeOF0GeNg",
							Tags:          []string{"data"},
							Timestamp:     1697091629,
						},
					},
				},
				Status:    true,
				Timestamp: 1697092032,
			},
			wantError: require.NoError,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			instance, err := worker.NewWorker(testcase.arguments.config)
			require.NoError(t, err)

			matched, err := instance.Match(ctx, testcase.arguments.task)
			testcase.wantError(t, err)
			require.True(t, matched)

			feed, err := instance.Transform(ctx, testcase.arguments.task)
			testcase.wantError(t, err)

			//t.Log(string(lo.Must(json.MarshalIndent(feed, "", "\x20\x20"))))

			require.Equal(t, testcase.want, feed)
		})
	}
}
