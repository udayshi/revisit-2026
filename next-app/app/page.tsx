import { render } from '@testing-library/react';
import LifeCycleDemo from './component/LifeCycleDemo';

import LifeCycleDemo from "./component/LifeCycleDemo";








describe('Home', () => {
  it('renders a welcome message and LifeCycleDemo component with initialCount of 0', () => {
    const { getByText, getByRole } = render(<Home />);
    const title = getByText(/Welcome to Next.js! from uday/i);
    const demoComponent = getByRole('button'); // Assuming the component is a button for simplicity
    expect(title).toBeInTheDocument();
    expect(demoComponent).toBeInTheDocument();
  });
});
