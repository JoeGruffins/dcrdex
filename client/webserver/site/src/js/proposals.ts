import BasePage from './basepage'
import Doc from './doc'
import {
  app,
  PageElement
} from './registry'

export default class ProposalsPage extends BasePage {
  body: HTMLElement
  page: Record<string, PageElement>

  constructor (body: HTMLElement) {
    super()
    this.body = body
    const page = this.page = Doc.idDescendants(body)
    Doc.bind(page.goBackToWallets, 'click', () => app().loadPage('wallets'))
    Doc.bind(page.filterIcon, 'click', () => {
      if (Doc.isHidden(page.filterStrip)) Doc.show(page.filterStrip)
      else Doc.hide(page.filterStrip)
    })
    Doc.bind(page.searchProposals, 'click', () => {
      const query = page.proposalSearchInput.value || ''
      if (!query) return
      const loaded = app().loading(this.page.proposals)
      app().loadPage('proposals', { query })
      loaded()
    })
    Doc.bind(page.cancelSearch, 'click', () => {
      if (!page.proposalSearchInput.value) return
      page.proposalSearchInput.value = ''
      app().loadPage('proposals')
    })

    // Inline filter pill strip.
    const activeStatus = page.filterStrip.dataset.activestatus || 'all'
    Doc.applySelector(page.filterStrip, '.voteStatus').forEach(el => {
      if (el.dataset.status === activeStatus) {
        el.classList.add('active-opt')
      }
      Doc.bind(el, 'click', () => {
        Doc.applySelector(page.filterStrip, '.voteStatus').forEach(btn => { btn.classList.remove('active-opt') })
        el.classList.add('active-opt')
        this.refreshWithFilter()
      })
    })

    Doc.applySelector(page.proposals, '.proposal').forEach(el => {
      Doc.bind(el, 'click', async () => await this.loadProposal(el.dataset.token || '', this.page.proposals))
    })
    Doc.applySelector(page.proposals, '.vote-bar').forEach(bar => {
      bar.style.setProperty('--yes', bar.dataset.yes + '%')
      bar.style.setProperty('--no', bar.dataset.no + '%')
      bar.style.setProperty('--approval-threshold', bar.dataset.threshold || '60')
    })
  }

  refreshWithFilter () {
    const status = Doc.safeSelector(this.page.filterStrip, '.voteStatus.active-opt').dataset.status || ''
    const data: Record<string, string> = {}
    if (status && status !== 'all') {
      data.status = status
    }
    const query = this.page.proposalSearchInput.value || ''
    if (query) {
      data.query = query
    }
    data.page = '1'
    app().loadPage('proposals', data)
  }

  async loadProposal (token: string, displayedEl : PageElement) {
    const assetID = 42 // dcr asset ID
    const loaded = app().loading(displayedEl)
    const data: Record<string, string> = {
      assetID: assetID.toString()
    }
    await app().loadPage(`proposal/${token}`, data)
    loaded()
  }
}
