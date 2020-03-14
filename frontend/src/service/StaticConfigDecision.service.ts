import { IFeatureDecisionService } from './FeatureDecision.service';

export class StaticConfigDecisionService implements IFeatureDecisionService {
  includeViewChangeLogButton(): boolean {
    return false;
  }
  includeSearchBar(): boolean {
    return false;
  }
  includeGithubSignButton(): boolean {
    return false;
  }
  includeGoogleSignButton(): boolean {
    return true;
  }
  includeFacebookSignButton(): boolean {
    return true;
  }
}
