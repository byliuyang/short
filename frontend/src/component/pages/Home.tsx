import React, { Component } from 'react';
import './Home.scss';

import { Header } from './shared/Header';
import { Url } from '../../entity/Url';
import { Footer } from './shared/Footer';
import { SignInModal } from './shared/sign-in/SignInModal';
import { Modal } from '../ui/Modal';
import { ExtPromo } from './shared/promos/ExtPromo';
import { CaptchaService } from '../../service/Captcha.service';
import { validateLongLinkFormat } from '../../validators/LongLink.validator';
import { validateCustomAliasFormat } from '../../validators/CustomAlias.validator';
import { Location } from 'history';
import { AuthService } from '../../service/Auth.service';
import { IBrowserExtensionService } from '../../service/extensionService/BrowserExtension.service';
import { VersionService } from '../../service/Version.service';
import { QrCodeService } from '../../service/QrCode.service';
import { UIFactory } from '../UIFactory';
import { IAppState } from '../../state/reducers';
import { Store } from 'redux';
import {
  clearError,
  raiseCreateShortLinkError,
  raiseInputError,
  updateAlias,
  updateCreatedUrl,
  updateLongLink
} from '../../state/actions';
import { ErrorService } from '../../service/Error.service';
import { IErr } from '../../entity/Err';
import { UrlService } from '../../service/Url.service';
import { SearchService } from '../../service/Search.service';
import { Update } from '../../entity/Update';
import { ChangeLogModal } from '../ui/ChangeLogModal';
import { CreateShortLinkSection } from './shared/CreateShortLinkSection';

interface Props {
  uiFactory: UIFactory;
  urlService: UrlService;
  authService: AuthService;
  extensionService: IBrowserExtensionService;
  versionService: VersionService;
  qrCodeService: QrCodeService;
  captchaService: CaptchaService;
  searchService: SearchService;
  errorService: ErrorService;
  store: Store<IAppState>;
  location: Location;
}

interface State {
  isUserSignedIn?: boolean;
  shouldShowPromo?: boolean;
  longLink?: string;
  alias?: string;
  shortLink?: string;
  createdUrl?: Url;
  qrCodeUrl?: string;
  err?: IErr;
  inputErr?: string;
  autoCompleteSuggestions?: Array<Url>;
  changeLog?: Array<Update>;
}

export class Home extends Component<Props, State> {
  errModal = React.createRef<Modal>();
  signInModal = React.createRef<SignInModal>();
  createShortLinkSection = React.createRef<CreateShortLinkSection>();
  changeLogModalRef = React.createRef<ChangeLogModal>();

  constructor(props: Props) {
    super(props);
    this.state = {
      changeLog: []
    };
  }

  componentDidMount(): void {
    this.setPromoDisplayStatus();

    this.props.authService.cacheAuthToken(this.props.location.search);
    if (!this.props.authService.isSignedIn()) {
      this.setState({
        isUserSignedIn: false
      });
      this.showSignInModal();
      return;
    }
    this.setState({
      isUserSignedIn: true
    });
    this.handleStateChange();
    this.autoFillLongLink();
  }

  async setPromoDisplayStatus() {
    var shouldShowPromo =
      this.props.extensionService.isSupported() &&
      !(await this.props.extensionService.isInstalled());
    this.setState({ shouldShowPromo: shouldShowPromo });
  }

  autoFillLongLink() {
    const longLink = this.getLongLinkFromQueryParams();
    if (validateLongLinkFormat(longLink) != null) {
      return;
    }
    this.props.store.dispatch(updateLongLink(longLink));
    this.createShortLinkSection.current!.focusShortLinkTextField();
  }

  handleStateChange() {
    this.props.store.subscribe(async () => {
      const state = this.props.store.getState();

      const newState: State = {
        longLink: state.editingUrl.originalUrl,
        alias: state.editingUrl.alias,
        err: state.err,
        createdUrl: state.createdUrl,
        inputErr: state.inputErr
      };

      if (state.createdUrl && state.createdUrl.alias) {
        const shortLink = this.props.urlService.aliasToFrontendLink(
          state.createdUrl.alias!
        );
        newState.shortLink = shortLink;
        newState.qrCodeUrl = await this.props.qrCodeService.newQrCode(
          shortLink
        );
      }

      if (newState.err) {
        console.log(newState.err);
        this.showError(newState.err);
      }
      this.setState(newState);
    });
  }

  showSignInModal() {
    if (!this.signInModal.current) {
      return;
    }
    this.signInModal.current.open();
  }

  requestSignIn = () => {
    this.setState({
      isUserSignedIn: false
    });
    this.props.authService.signOut();
    this.showSignInModal();
  };

  handleSearchBarInputChange = async (alias: String) => {
    const autoCompleteSuggestions = await this.props.searchService.getAutoCompleteSuggestions(
      alias
    );
    this.setState({
      autoCompleteSuggestions
    });
  };

  handleSignOutButtonClick = () => {
    this.requestSignIn();
  };

  handlerLongLinkChange = (newLongLink: string) => {
    this.props.store.dispatch(updateLongLink(newLongLink));
  };

  handleAliasChange = (newAlias: string) => {
    this.props.store.dispatch(updateAlias(newAlias));
  };

  handleOnErrModalCloseClick = () => {
    this.errModal.current!.close();
    this.props.store.dispatch(clearError());
  };

  handlerLongLinkTextFieldBlur = () => {
    let longLink = this.props.store.getState().editingUrl.originalUrl;
    let err = validateLongLinkFormat(longLink);
    this.props.store.dispatch(raiseInputError(err));
  };

  handlerCustomAliasTextFieldBlur = () => {
    const alias = this.props.store.getState().editingUrl.alias;
    const err = validateCustomAliasFormat(alias);
    this.props.store.dispatch(raiseInputError(err));
  };

  handleCreateShortLinkClick = () => {
    const editingUrl = this.props.store.getState().editingUrl;
    this.props.urlService
      .createShortLink(editingUrl)
      .then((createdUrl: Url) =>
        this.props.store.dispatch(updateCreatedUrl(createdUrl))
      )
      .catch(({ authorizationErr, createShortLinkErr }) => {
        if (authorizationErr) {
          this.requestSignIn();
          return;
        }
        this.props.store.dispatch(
          raiseCreateShortLinkError(createShortLinkErr)
        );
      });
  };

  getLongLinkFromQueryParams(): string {
    let urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('long_link')!;
  }

  showError(error?: IErr) {
    if (!error) {
      return;
    }
    this.errModal.current!.open();
  }

  handleShowChangeLogBtnClick = () => {
    if (this.changeLogModalRef.current) {
      this.changeLogModalRef.current.open();
    }
  };

  render = () => {
    return (
      <div className="home">
        {this.state.shouldShowPromo && <ExtPromo />}
        <Header
          uiFactory={this.props.uiFactory}
          onSearchBarInputChange={this.handleSearchBarInputChange}
          autoCompleteSuggestions={this.state.autoCompleteSuggestions}
          shouldShowSignOutButton={this.state.isUserSignedIn}
          onSignOutButtonClick={this.handleSignOutButtonClick}
        />
        <div className={'main'}>
          <CreateShortLinkSection
            longLinkText={this.state.longLink}
            alias={this.state.alias}
            shortLink={this.state.shortLink}
            inputErr={this.state.inputErr}
            createdUrl={this.state.createdUrl}
            qrCodeUrl={this.state.qrCodeUrl}
            onLongLinkTextFieldBlur={this.handlerLongLinkTextFieldBlur}
            onLongLinkTextFieldChange={this.handlerLongLinkChange}
            onShortLinkTextFieldBlur={this.handlerCustomAliasTextFieldBlur}
            onShortLinkTextFieldChange={this.handleAliasChange}
            onCreateShortLinkButtonClick={this.handleCreateShortLinkClick}
          />
        </div>
        <Footer
          uiFactory={this.props.uiFactory}
          onShowChangeLogBtnClick={this.handleShowChangeLogBtnClick}
          authorName={'Harry'}
          authorPortfolio={'https://github.com/byliuyang'}
          version={this.props.versionService.getAppVersion()}
        />
        <ChangeLogModal
          ref={this.changeLogModalRef}
          changeLog={this.state.changeLog}
          defaultVisibleLogs={3}
        />

        <SignInModal ref={this.signInModal} uiFactory={this.props.uiFactory} />
        <Modal canClose={true} ref={this.errModal}>
          {this.state.err ? (
            <div className={'err'}>
              <i
                className={'material-icons close'}
                title={'close'}
                onClick={this.handleOnErrModalCloseClick}
              >
                close
              </i>
              <div className={'title'}>{this.state.err.name}</div>
              <div className={'description'}>{this.state.err.description}</div>
            </div>
          ) : (
            false
          )}
        </Modal>
      </div>
    );
  };
}
